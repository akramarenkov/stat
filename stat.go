package stat

import (
	"io"
	"strconv"

	"github.com/akramarenkov/safe"
	"github.com/akramarenkov/safe/intspec"
	"github.com/akramarenkov/safe/is"
	"github.com/akramarenkov/span"
	"github.com/pterm/pterm"
	"golang.org/x/exp/constraints"
)

// Statistics.
type Stat[Type constraints.Integer] struct {
	negInf    Item[Type]
	posInf    Item[Type]
	predictor Predictor[Type]
	items     []Item[Type]
}

// Creates an instance of statistics for the specified spans of values.
//
// Spans of values must be sorted from smallest to largest and must not intersect.
//
// Prediction function may not be specified, but then the value's correspondence to
// the span will be determined by searching the list of spans, which is slower.
func New[Type constraints.Integer](spans []span.Span[Type], predictor Predictor[Type]) (*Stat[Type], error) {
	if len(spans) == 0 {
		return nil, ErrSpansEmpty
	}

	st := &Stat[Type]{
		items:     createItems(spans),
		predictor: predictor,
	}

	st.initInf()

	return st, nil
}

func createItems[Type constraints.Integer](spans []span.Span[Type]) []Item[Type] {
	items := make([]Item[Type], len(spans))

	for id, span := range spans {
		items[id].Span = span
	}

	return items
}

func (st *Stat[Type]) initInf() {
	minimum, maximum := intspec.Range[Type]()

	if minimum < st.items[0].Span.Begin {
		negative := span.Span[Type]{
			Begin: minimum,
			End:   st.items[0].Span.Begin - 1,
		}

		st.negInf.Span = negative
	}

	if maximum > st.items[len(st.items)-1].Span.End {
		positive := span.Span[Type]{
			Begin: st.items[len(st.items)-1].Span.End + 1,
			End:   maximum,
		}

		st.posInf.Span = positive
	}
}

// Increases the quantity of occurrences of a value within the specified spans.
func (st *Stat[Type]) Inc(value Type) {
	if value < st.items[0].Span.Begin {
		// Integer overflow is possible here and below, but it will take a long time
		// and this case cannot be tested
		st.negInf.Quantity++
		return
	}

	if value > st.items[len(st.items)-1].Span.End {
		st.posInf.Quantity++
		return
	}

	if st.predictor != nil {
		st.items[st.predictor(value)].Quantity++
		return
	}
}

// Returns a list of statistics items.
func (st *Stat[Type]) Items() []Item[Type] {
	items := make([]Item[Type], len(st.items)+infinitiesQuantity)

	items[0] = st.negInf

	copy(items[1:], st.items)

	items[len(items)-1] = st.posInf

	if st.negInf.Quantity == 0 {
		items = items[1:]
	}

	if st.posInf.Quantity == 0 {
		items = items[:len(items)-1]
	}

	return items
}

// Writes statistics as a bar chart to the specified writers.
func (st *Stat[Type]) Graph(writers ...io.Writer) error {
	bars := make([]pterm.Bar, 0, len(st.items))

	style := &pterm.Style{
		pterm.BgDefault,
		pterm.FgDefault,
	}

	if st.negInf.Quantity != 0 {
		value, err := safe.IToI[int](st.negInf.Quantity)
		if err != nil {
			return err
		}

		bar := pterm.Bar{
			Label:      "[-Inf:" + fmtInt(st.negInf.Span.End) + "]",
			Value:      value,
			Style:      style,
			LabelStyle: style,
		}

		bars = append(bars, bar)
	}

	for _, item := range st.items {
		value, err := safe.IToI[int](item.Quantity)
		if err != nil {
			return err
		}

		bar := pterm.Bar{
			Label:      "[" + fmtInt(item.Span.Begin) + ":" + fmtInt(item.Span.End) + "]",
			Value:      value,
			Style:      style,
			LabelStyle: style,
		}

		bars = append(bars, bar)
	}

	if st.posInf.Quantity != 0 {
		value, err := safe.IToI[int](st.posInf.Quantity)
		if err != nil {
			return err
		}

		bar := pterm.Bar{
			Label:      "[" + fmtInt(st.posInf.Span.Begin) + ":+Inf" + "]",
			Value:      value,
			Style:      style,
			LabelStyle: style,
		}

		bars = append(bars, bar)
	}

	chart := pterm.DefaultBarChart.WithBars(bars).WithShowValue()

	if len(writers) == 0 {
		// In the library version used, this function actually never returns errors
		_ = chart.Render()
	}

	for _, writer := range writers {
		_ = chart.WithWriter(writer).Render()
	}

	return nil
}

func fmtInt[Type constraints.Integer](number Type) string {
	if is.Signed[Type]() {
		return strconv.FormatInt(int64(number), decimalBase)
	}

	return strconv.FormatUint(uint64(number), decimalBase)
}

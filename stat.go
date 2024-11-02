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

	lower := st.items[st.lower()]
	upper := st.items[st.upper()]

	if minimum < lower.Span.Begin {
		negInf := span.Span[Type]{
			Begin: minimum,
			End:   lower.Span.Begin - 1,
		}

		st.negInf.Span = negInf
	}

	if maximum > upper.Span.End {
		posInf := span.Span[Type]{
			Begin: upper.Span.End + 1,
			End:   maximum,
		}

		st.posInf.Span = posInf
	}
}

// Increases the quantity of occurrences of a value within the specified spans.
func (st *Stat[Type]) Inc(value Type) {
	if value < st.items[st.lower()].Span.Begin {
		// Integer overflow is possible here and below, but it will take a long time
		// and this case cannot be tested
		st.negInf.Quantity++
		return
	}

	if value > st.items[st.upper()].Span.End {
		st.posInf.Quantity++
		return
	}

	if st.predictor != nil {
		st.items[st.predictor(value)].Quantity++
		return
	}
}

func (st *Stat[Type]) lower() int {
	return 0
}

func (st *Stat[Type]) upper() int {
	return len(st.items) - 1
}

// Returns a list of statistics items.
func (st *Stat[Type]) Items() []Item[Type] {
	items := make([]Item[Type], len(st.items)+infinitiesQuantity)

	baseOffest := 1
	negInfID := 0
	posInfID := len(items) - 1

	copy(items[baseOffest:], st.items)

	items[negInfID] = st.negInf
	items[posInfID] = st.posInf

	// Order of trimming is important because the offset from the top depends on the
	// length of the slice
	if st.posInf.Quantity == 0 {
		items = items[:posInfID]
	}

	if st.negInf.Quantity == 0 {
		items = items[baseOffest:]
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

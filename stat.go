package stat

import (
	"io"
	"os"
	"slices"
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
	items     []Item[Type]
	missed    Item[Type]
	negInf    Item[Type]
	posInf    Item[Type]
	predictor Predictor[Type]
}

// Creates an instance of statistics for the specified spans of values.
//
// Spans sequence must be increasing and sorted. Spans must not intersect.
//
// Prediction function may not be specified, but then the value's correspondence to
// the span will be determined by searching the list of spans, which is slower.
func New[Type constraints.Integer](spans []span.Span[Type], predictor Predictor[Type]) (*Stat[Type], error) {
	if len(spans) == 0 {
		return nil, ErrSpansListEmpty
	}

	if err := span.IsNotIntersect(spans); err != nil {
		return nil, err
	}

	if err := span.IsNonDecreasing(spans); err != nil {
		return nil, err
	}

	if !slices.IsSortedFunc(spans, span.CompareInc) {
		return nil, ErrSpansSequenceUnsorted
	}

	st := &Stat[Type]{
		items:     createItems(spans),
		predictor: predictor,
	}

	st.prepare()

	return st, nil
}

func createItems[Type constraints.Integer](spans []span.Span[Type]) []Item[Type] {
	items := make([]Item[Type], len(spans))

	for id, span := range spans {
		items[id].Span = span
		items[id].Kind = ItemKindRegular
	}

	return items
}

func (st *Stat[Type]) prepare() {
	st.missed.Kind = ItemKindMissed
	st.negInf.Kind = ItemKindNegInf
	st.posInf.Kind = ItemKindPosInf

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

// Increases the quantity of occurrences of the specified value.
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

	target := Item[Type]{
		Span: span.Span[Type]{Begin: value, End: value},
	}

	if id, found := slices.BinarySearchFunc(st.items, target, search); found {
		st.items[id].Quantity++
		return
	}

	st.missed.Quantity++
}

func (st *Stat[Type]) lower() int {
	return 0
}

func (st *Stat[Type]) upper() int {
	return len(st.items) - 1
}

// Returns a list of statistics items.
func (st *Stat[Type]) Items() []Item[Type] {
	items := make([]Item[Type], 0, len(st.items)+specialItemsQuantity)

	if st.missed.Quantity != 0 {
		items = append(items, st.missed)
	}

	if st.negInf.Quantity != 0 {
		items = append(items, st.negInf)
	}

	items = append(items, st.items...)

	if st.posInf.Quantity != 0 {
		items = append(items, st.posInf)
	}

	return items
}

// Writes statistics as a bar chart to the specified writers.
//
// If no writer is specified, the bar chart will be written to standard output.
func (st *Stat[Type]) Graph(writers ...io.Writer) error {
	if len(writers) == 0 {
		if err := st.graph(os.Stdout); err != nil {
			return err
		}

		return nil
	}

	for _, writer := range writers {
		if err := st.graph(writer); err != nil {
			return err
		}
	}

	return nil
}

func (st *Stat[Type]) graph(writer io.Writer) error {
	bars := make([]pterm.Bar, 0, len(st.items)+specialItemsQuantity)

	style := &pterm.Style{
		pterm.BgDefault,
		pterm.FgDefault,
	}

	if st.missed.Quantity != 0 {
		value, err := safe.IToI[int](st.missed.Quantity)
		if err != nil {
			return err
		}

		bar := pterm.Bar{
			Label:      "[" + st.missed.Kind.String() + "]",
			Value:      value,
			Style:      style,
			LabelStyle: style,
		}

		bars = append(bars, bar)
	}

	if st.negInf.Quantity != 0 {
		value, err := safe.IToI[int](st.negInf.Quantity)
		if err != nil {
			return err
		}

		bar := pterm.Bar{
			Label:      "[" + st.negInf.Kind.String() + ":" + fmtInt(st.negInf.Span.End) + "]",
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
			Label:      "[" + fmtInt(st.posInf.Span.Begin) + ":" + st.posInf.Kind.String() + "]",
			Value:      value,
			Style:      style,
			LabelStyle: style,
		}

		bars = append(bars, bar)
	}

	chart := pterm.DefaultBarChart.WithBars(bars).WithShowValue()

	return chart.WithWriter(writer).Render()
}

func fmtInt[Type constraints.Integer](number Type) string {
	if is.Signed[Type]() {
		return strconv.FormatInt(int64(number), decimalBase)
	}

	return strconv.FormatUint(uint64(number), decimalBase)
}

package stat

import (
	"io"
	"slices"
	"strconv"

	"github.com/akramarenkov/safe"
	"github.com/akramarenkov/safe/intspec"
	"github.com/akramarenkov/safe/is"
	"github.com/akramarenkov/span"
	"github.com/pterm/pterm"
	"golang.org/x/exp/constraints"
)

// Item of statistics.
type Item[Type constraints.Integer] struct {
	// Quantity of occurrences of a value from a Span
	Quantity uint64
	// Span of values ​​for which the Quantity of occurrences is collected
	Span span.Span[Type]
}

// Statistics.
type Stat[Type constraints.Integer] struct {
	items           []Item[Type]
	negInfIsPresent bool
	posInfIsPresent bool
}

// Creates an instance of statistics for the specified spans of values.
//
// Spans of values must be sorted from smallest to largest.
func New[Type constraints.Integer](spans []span.Span[Type]) *Stat[Type] {
	st := &Stat[Type]{
		items: make([]Item[Type], len(spans)+infinitiesQuantity),
	}

	if len(spans) == 0 {
		return st
	}

	lowerID := 0
	upperID := len(spans) - 1

	minimum, maximum := intspec.Range[Type]()

	if minimum < spans[lowerID].Begin {
		negative := span.Span[Type]{
			Begin: minimum,
			End:   spans[lowerID].Begin - 1,
		}

		st.items[st.negInfID()].Span = negative
		st.negInfIsPresent = true
	}

	if maximum > spans[upperID].End {
		positive := span.Span[Type]{
			Begin: spans[upperID].End + 1,
			End:   maximum,
		}

		st.items[st.posInfID()].Span = positive
		st.posInfIsPresent = true
	}

	for id, span := range spans {
		st.items[id+st.baseOffsetInt()].Span = span
	}

	return st
}

// Increases the quantity of occurrences of a value lesser than the smallest value in
// the specified spans.
func (st *Stat[Type]) IncNegInf() {
	if !st.negInfIsPresent {
		panic(ErrNegInfNotPresent)
	}

	st.inc(st.negInfID())
}

// Increases the quantity of occurrences of a value greater than the largest value in
// the specified spans.
func (st *Stat[Type]) IncPosInf() {
	if !st.posInfIsPresent {
		panic(ErrPosInfNotPresent)
	}

	st.inc(st.posInfID())
}

// Increases the quantity of occurrences of a value within the specified spans.
//
// It is necessary to specify the index of the span in which the value hits.
func (st *Stat[Type]) Inc(id uint64) {
	// Integer overflow is possible here, but this will require as much RAM as does
	// not exist yet in one machine and this case cannot be tested
	id += st.baseOffset()

	if id >= st.posInfID() {
		panic(ErrInvalidIndex)
	}

	st.inc(id)
}

func (st *Stat[Type]) inc(id uint64) {
	// Integer overflow is possible here, but it will take a long time and this case
	// cannot be tested
	st.items[id].Quantity++
}

func (st *Stat[Type]) negInfID() uint64 {
	return 0
}

func (st *Stat[Type]) posInfID() uint64 {
	return uint64(len(st.items)) - 1
}

func (st *Stat[Type]) baseOffset() uint64 {
	return 1
}

func (st *Stat[Type]) baseOffsetInt() int {
	return 1
}

func (st *Stat[Type]) base() []Item[Type] {
	return st.items[st.baseOffset():st.posInfID()]
}

// Returns a list of statistics items.
func (st *Stat[Type]) Items() []Item[Type] {
	items := st.items

	if !st.posInfIsPresent {
		items = items[:st.posInfID()]
	}

	if !st.negInfIsPresent {
		items = items[st.baseOffset():]
	}

	return slices.Clone(items)
}

// Writes statistics as a bar chart to the specified writers.
func (st *Stat[Type]) Graph(writers ...io.Writer) error {
	bars := make([]pterm.Bar, 0, len(st.items))

	style := &pterm.Style{
		pterm.BgDefault,
		pterm.FgDefault,
	}

	if negInf := st.items[st.negInfID()]; negInf.Quantity != 0 {
		value, err := safe.IToI[int](negInf.Quantity)
		if err != nil {
			return err
		}

		bar := pterm.Bar{
			Label:      "[-Inf:" + fmtInt(negInf.Span.End) + "]",
			Value:      value,
			Style:      style,
			LabelStyle: style,
		}

		bars = append(bars, bar)
	}

	for _, item := range st.base() {
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

	if posInf := st.items[st.posInfID()]; posInf.Quantity != 0 {
		value, err := safe.IToI[int](posInf.Quantity)
		if err != nil {
			return err
		}

		bar := pterm.Bar{
			Label:      "[" + fmtInt(posInf.Span.Begin) + ":+Inf" + "]",
			Value:      value,
			Style:      style,
			LabelStyle: style,
		}

		bars = append(bars, bar)
	}

	chart := pterm.DefaultBarChart.WithBars(bars).WithShowValue()

	for _, writer := range writers {
		chart = chart.WithWriter(writer)

		// The version used actually does not return errors
		_ = chart.Render()
	}

	return nil
}

func fmtInt[Type constraints.Integer](number Type) string {
	if is.Signed[Type]() {
		return strconv.FormatInt(int64(number), decimalBase)
	}

	return strconv.FormatUint(uint64(number), decimalBase)
}

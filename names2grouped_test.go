package names2grouped_test

import (
	"iter"
	"slices"
	"testing"

	ng "github.com/takanoriyanagitani/go-names2grouped"
)

func TestNamesToGrouped(t *testing.T) {
	t.Parallel()

	t.Run("IsSameGroup", func(t *testing.T) {
		t.Parallel()

		t.Run("FirstStringCheck", func(t *testing.T) {
			t.Parallel()

			t.Run("empty", func(t *testing.T) {
				t.Parallel()

				var empty []string
				var names iter.Seq[string] = slices.Values(empty)

				var chk1st ng.IsSameGroup = ng.FirstStringCheck
				var grouped ng.GroupedNames = chk1st.
					NamesToGrouped(ng.Names(names))
				var collected [][]string = grouped.Collect()
				if 0 != len(collected) {
					t.Fatalf("must be empty\n")
				}
			})

			t.Run("single", func(t *testing.T) {
				t.Parallel()

				var single []string = []string{"helo"}
				var names iter.Seq[string] = slices.Values(single)

				var chk1st ng.IsSameGroup = ng.FirstStringCheck
				var grouped ng.GroupedNames = chk1st.
					NamesToGrouped(ng.Names(names))
				var collected [][]string = grouped.Collect()
				if 1 != len(collected) {
					t.Fatalf("unexpected len: %v\n", len(collected))
				}
			})

			t.Run("single-group multi-names", func(t *testing.T) {
				t.Parallel()

				var multi []string = []string{"helo", "hell"}
				var names iter.Seq[string] = slices.Values(multi)

				var chk1st ng.IsSameGroup = ng.FirstStringCheck
				var grouped ng.GroupedNames = chk1st.
					NamesToGrouped(ng.Names(names))
				var collected [][]string = grouped.Collect()
				if 1 != len(collected) {
					t.Fatalf("unexpected len: %v\n", len(collected))
				}
			})

			t.Run("multi-group multi-names", func(t *testing.T) {
				t.Parallel()

				var multi []string = []string{
					"fuji", "fuji9",
					"hello", "heli",
					"tokyo", "takao", "tree",
				}
				var names iter.Seq[string] = slices.Values(multi)

				var chk1st ng.IsSameGroup = ng.FirstStringCheck
				var grouped ng.GroupedNames = chk1st.
					NamesToGrouped(ng.Names(names))
				var collected [][]string = grouped.Collect()
				if 3 != len(collected) {
					t.Fatalf("unexpected len: %v\n", len(collected))
				}

				tables := []struct {
					ix1      int
					ix2      int
					expected string
				}{
					{ix1: 0, ix2: 0, expected: "fuji"},
					{ix1: 0, ix2: 1, expected: "fuji9"},
					{ix1: 1, ix2: 0, expected: "hello"},
					{ix1: 1, ix2: 1, expected: "heli"},
					{ix1: 2, ix2: 0, expected: "tokyo"},
					{ix1: 2, ix2: 1, expected: "takao"},
					{ix1: 2, ix2: 2, expected: "tree"},
				}

				for _, tab := range tables {
					if collected[tab.ix1][tab.ix2] != tab.expected {
						t.Fatalf("unexpected value: %v\n", collected[tab.ix1])
					}
				}
			})
		})

		t.Run("FirstBytesCheckNew", func(t *testing.T) {
			t.Parallel()

			t.Run("empty", func(t *testing.T) {
				t.Parallel()

				var empty []string
				var names iter.Seq[string] = slices.Values(empty)

				var chk ng.IsSameGroup = ng.FirstBytesCheckNew(4)
				var grouped ng.GroupedNames = chk.
					NamesToGrouped(ng.Names(names))
				var collected [][]string = grouped.Collect()
				if 0 != len(collected) {
					t.Fatalf("must be empty\n")
				}
			})

			t.Run("single, short name", func(t *testing.T) {
				t.Parallel()

				var short []string = []string{"he"}
				var names iter.Seq[string] = slices.Values(short)

				var chk ng.IsSameGroup = ng.FirstBytesCheckNew(4)
				var grouped ng.GroupedNames = chk.
					NamesToGrouped(ng.Names(names))
				var collected [][]string = grouped.Collect()
				if 1 != len(collected) {
					t.Fatalf("unexpected len: %v\n", len(collected))
				}
			})

			t.Run("single group", func(t *testing.T) {
				t.Parallel()

				var single []string = []string{
					"helo, wrld",
					"helo, helo",
					"helo, world",
				}
				var names iter.Seq[string] = slices.Values(single)

				var chk ng.IsSameGroup = ng.FirstBytesCheckNew(4)
				var grouped ng.GroupedNames = chk.
					NamesToGrouped(ng.Names(names))
				var collected [][]string = grouped.Collect()
				if 1 != len(collected) {
					t.Fatalf("unexpected len: %v\n", len(collected))
				}

				var results []string = collected[0]
				var l int = len(results)
				if l != len(single) {
					t.Fatalf("unexpected len: %v\n", l)
				}

				for idx := range l {
					var si string = single[idx]
					var ri string = results[idx]
					if si != ri {
						t.Fatalf("unexpected value got: %s != %s\n", si, ri)
					}
				}
			})

			t.Run("multi group", func(t *testing.T) {
				t.Parallel()

				var multi []string = []string{
					"unsorted item 1",
					"unsorted item 2",
					"helo, wrld",
					"helo, helo",
					"helo, world",
					"unsorted item 3",
					"unsorted item 4",
				}
				var names iter.Seq[string] = slices.Values(multi)

				var chk ng.IsSameGroup = ng.FirstBytesCheckNew(4)
				var grouped ng.GroupedNames = chk.
					NamesToGrouped(ng.Names(names))
				var collected [][]string = grouped.Collect()
				if 3 != len(collected) {
					t.Fatalf("unexpected len: %v\n", len(collected))
				}

				rows := []struct {
					i0       int
					i1       int
					expected string
				}{
					{i0: 0, i1: 0, expected: "unsorted item 1"},
					{i0: 0, i1: 1, expected: "unsorted item 2"},

					{i0: 1, i1: 0, expected: "helo, wrld"},
					{i0: 1, i1: 1, expected: "helo, helo"},
					{i0: 1, i1: 2, expected: "helo, world"},

					{i0: 2, i1: 0, expected: "unsorted item 3"},
					{i0: 2, i1: 1, expected: "unsorted item 4"},
				}

				for _, row := range rows {
					var i0 int = row.i0
					var i1 int = row.i1
					var got string = collected[i0][i1]
					if got != row.expected {
						t.Fatalf(
							"unexpected value got. %s != %s\n",
							got,
							row.expected,
						)
					}
				}
			})
		})
	})
}

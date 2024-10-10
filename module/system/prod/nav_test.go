package prod

import "testing"

func TestRequireN(t *testing.T) {
	var menus []Menu = []Menu{
		{
			Name: "1",
			Path: "/",
		},
		{
			Name: "2",
			Path: "/2",
			Children: []Menu{
				{
					Name:    "2.1",
					Path:    "/2.1",
					Require: true,
				},
				{
					Name: "2.2",
					Path: "/2.2",
				},
			},
		},
		{
			Name: "3",
			Path: "/3",
			Children: []Menu{
				{
					Name: "3.1",
					Path: "/3.1",
					Children: []Menu{
						{
							Name:    "3.1.1",
							Path:    "/3.1.1",
							Require: true,
						},
					},
				},
			},
		},
	}
	requires := RequireNav(menus)
	if len(requires) != 2 {
		t.Errorf("require nav 1 faile\n")
		return
	}
	if len(requires[0].Children) != 1 {
		t.Errorf("require nav 2 faile\n")
		return
	}
	if len(requires[1].Children[0].Children) != 1 {
		t.Errorf("require nav 3 faile\n")
		return
	}
}

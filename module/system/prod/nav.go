package prod

func RequireNav(menus []Menu) []Menu {
	var temps []Menu = []Menu{}
	for index := range menus {
		menu := menus[index]
		if len(menu.Children) != 0 {
			menu.Children = RequireNav(menu.Children)
		}
		if menu.Require || len(menu.Children) != 0 {
			temps = append(temps, menu)
		}
	}
	return temps
}

func GetRequireNav() []Menu {
	return RequireNav(gNavConf.Menus)
}

func GetNav() []Menu {
	return gNavConf.Menus
}

package menu

// MenuItem represents a menuitem contained in a menu
type MenuItem struct {
	// The unique identifier of this menu item
	ID string `json:"ID,omitempty"`
	// Label is what appears as the menu text
	Label string
	// Role is a predefined menu type
	Role Role `json:"Role,omitempty"`
	// Accelerator holds a representation of a key binding
	Accelerator *Accelerator `json:"Accelerator,omitempty"`
	// Type of MenuItem, EG: Checkbox, Text, Separator, Radio, Submenu
	Type Type
	// Disabled makes the item unselectable
	Disabled bool
	// Hidden ensures that the item is not shown in the menu
	Hidden bool
	// Checked indicates if the item is selected (used by Checkbox and Radio types only)
	Checked bool
	// Submenu contains a list of menu items that will be shown as a submenu
	SubMenu []*MenuItem `json:"SubMenu,omitempty"`

	// This holds the menu item's parent.
	parent *MenuItem
}

// Parent returns the parent of the menu item.
// If it is a top level menu then it returns nil.
func (m *MenuItem) Parent() *MenuItem {
	return m.parent
}

// Append will attempt to append the given menu item to
// this item's submenu items. If this menu item is not a
// submenu, then this method will not add the item and
// simply return false.
func (m *MenuItem) Append(item *MenuItem) bool {
	if !m.isSubMenu() {
		return false
	}
	item.parent = m
	m.SubMenu = append(m.SubMenu, item)
	return true
}

// Prepend will attempt to prepend the given menu item to
// this item's submenu items. If this menu item is not a
// submenu, then this method will not add the item and
// simply return false.
func (m *MenuItem) Prepend(item *MenuItem) bool {
	if !m.isSubMenu() {
		return false
	}
	item.parent = m
	m.SubMenu = append([]*MenuItem{item}, m.SubMenu...)
	return true
}

func (m *MenuItem) getByID(id string) *MenuItem {

	// If I have the ID return me!
	if m.ID == id {
		return m
	}

	// Check submenus
	for _, submenu := range m.SubMenu {
		result := submenu.getByID(id)
		if result != nil {
			return result
		}
	}

	return nil
}

func (m *MenuItem) removeByID(id string) bool {

	for index, item := range m.SubMenu {
		if item.ID == id {
			m.SubMenu = append(m.SubMenu[:index], m.SubMenu[index+1:]...)
			return true
		}
		if item.isSubMenu() {
			result := item.removeByID(id)
			if result == true {
				return result
			}
		}
	}
	return false
}

// InsertAfter attempts to add the given item after this item in the parent
// menu. If there is no parent menu (we are a top level menu) then false is
// returned
func (m *MenuItem) InsertAfter(item *MenuItem) bool {

	// We need to find my parent
	if m.parent == nil {
		return false
	}

	// Get my parent to insert the item
	return m.parent.insertNewItemAfterGivenItem(m, item)
}

// InsertBefore attempts to add the given item before this item in the parent
// menu. If there is no parent menu (we are a top level menu) then false is
// returned
func (m *MenuItem) InsertBefore(item *MenuItem) bool {

	// We need to find my parent
	if m.parent == nil {
		return false
	}

	// Get my parent to insert the item
	return m.parent.insertNewItemBeforeGivenItem(m, item)
}

// insertNewItemAfterGivenItem will insert the given item after the given target
// in this item's submenu. If we are not a submenu,
// then something bad has happened :/
func (m *MenuItem) insertNewItemAfterGivenItem(target *MenuItem,
	newItem *MenuItem) bool {

	if !m.isSubMenu() {
		return false
	}

	// Find the index of the target
	targetIndex := m.getItemIndex(target)
	if targetIndex == -1 {
		return false
	}

	// Insert element into slice
	return m.insertItemAtIndex(targetIndex+1, newItem)
}

// insertNewItemBeforeGivenItem will insert the given item before the given
// target in this item's submenu. If we are not a submenu, then something bad
// has happened :/
func (m *MenuItem) insertNewItemBeforeGivenItem(target *MenuItem,
	newItem *MenuItem) bool {

	if !m.isSubMenu() {
		return false
	}

	// Find the index of the target
	targetIndex := m.getItemIndex(target)
	if targetIndex == -1 {
		return false
	}

	// Insert element into slice
	return m.insertItemAtIndex(targetIndex, newItem)
}

func (m *MenuItem) isSubMenu() bool {
	return m.Type == SubmenuType
}

// getItemIndex returns the index of the given target relative to this menu
func (m *MenuItem) getItemIndex(target *MenuItem) int {

	// This should only be called on submenus
	if !m.isSubMenu() {
		return -1
	}

	// hunt down that bad boy
	for index, item := range m.SubMenu {
		if item == target {
			return index
		}
	}

	return -1
}

// insertItemAtIndex attempts to insert the given item into the submenu at
// the given index
// Credit: https://stackoverflow.com/a/61822301
func (m *MenuItem) insertItemAtIndex(index int, target *MenuItem) bool {

	// If index is OOB, return false
	if index > len(m.SubMenu) {
		return false
	}

	// Save parent reference
	target.parent = m

	// If index is last item, then just regular append
	if index == len(m.SubMenu) {
		m.SubMenu = append(m.SubMenu, target)
		return true
	}

	m.SubMenu = append(m.SubMenu[:index+1], m.SubMenu[index:]...)
	m.SubMenu[index] = target
	return true
}

// Text is a helper to create basic Text menu items
func Text(label string, id string, accelerator *Accelerator) *MenuItem {
	return &MenuItem{
		ID:          id,
		Label:       label,
		Type:        TextType,
		Accelerator: accelerator,
	}
}

// Separator provides a menu separator
func Separator() *MenuItem {
	return &MenuItem{
		Type: SeparatorType,
	}
}

// Radio is a helper to create basic Radio menu items with an accelerator
func Radio(label string, id string, selected bool, accelerator *Accelerator) *MenuItem {
	return &MenuItem{
		ID:          id,
		Label:       label,
		Type:        RadioType,
		Checked:     selected,
		Accelerator: accelerator,
	}
}

// Checkbox is a helper to create basic Checkbox menu items
func Checkbox(label string, id string, checked bool, accelerator *Accelerator) *MenuItem {
	return &MenuItem{
		ID:          id,
		Label:       label,
		Type:        CheckboxType,
		Checked:     checked,
		Accelerator: accelerator,
	}
}

// SubMenu is a helper to create Submenus
func SubMenu(label string, items []*MenuItem) *MenuItem {
	result := &MenuItem{
		Label:   label,
		SubMenu: items,
		Type:    SubmenuType,
	}

	// Fix up parent pointers
	for _, item := range items {
		item.parent = result
	}

	return result
}

// SubMenuWithID is a helper to create Submenus with an ID
func SubMenuWithID(label string, id string, items []*MenuItem) *MenuItem {
	result := &MenuItem{
		Label:   label,
		SubMenu: items,
		ID:      id,
		Type:    SubmenuType,
	}

	// Fix up parent pointers
	for _, item := range items {
		item.parent = result
	}

	return result
}
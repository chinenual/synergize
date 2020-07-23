/* adapted from https://stackoverflow.com/questions/32636750/how-to-add-a-right-click-menu-in-electron-that-has-inspect-element-option-like */

const remote = require('electron').remote;
const Menu = remote.Menu;
const MenuItem = remote.MenuItem;

let debug = {
    rightClickPosition: null,

    configDebugContextMenu: function () {
        const menu = new Menu()
        const menuItem = new MenuItem({
            label: 'Inspect Element',
            click: () => {
                remote.getCurrentWindow().inspectElement(debug.rightClickPosition.x, debug.rightClickPosition.y)
            }
        })

        menu.append(menuItem)
        window.addEventListener('contextmenu', (e) => {
            e.preventDefault()
            debug.rightClickPosition = { x: e.x, y: e.y }
            menu.popup(remote.getCurrentWindow())
        }, false)
    }
};

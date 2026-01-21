// Adapted from
// https://github.com/martinacarter/vanilla-js-tabs/blob/main/tabs.html Uses
// bootstrap classnames to minimize HTML changes

export function tabsInit() {


        const tabsContainer = document.querySelector('.tab-content');
        const tabsList = document.querySelector('.nav-pills');
        const tabButtons = tabsList.querySelectorAll('.nav-link');
        const tabPanels = tabsContainer.querySelectorAll('.tab-pane');

        // guidelines of how accessible tab system should work

        tabPanels.forEach((panel) => {
                panel.setAttribute('tabindex', '0');
        });
        
        tabButtons.forEach((tab) => {
                tab.addEventListener('click', function (e) {
                        const clickedTab = e.target;

                        e.preventDefault();
                        setActiveTab(clickedTab);
                });
        });

        function setActiveTab(clickedTab) {
                // Get the active tab and associated active panel
                // We do this by grabbing href of the clicked tab
                // because that it was we used for the ID's of the associated panel
                const activeTab = clickedTab.getAttribute('href');
                console.log('Switch to tab ', activeTab);
                const activePanel = document.querySelector(activeTab);

                // RESET TABS: go through and make sure no tab is selected
                tabButtons.forEach((tab) => {
                        tab.setAttribute('aria-selected', false);
                        tab.setAttribute('tabindex', '-1');
                });
                // we'll also go through all the panels and add hidden to all panels
                tabPanels.forEach((tabPanel) => {
                        tabPanel.setAttribute('hidden', true);
                });

                // next we'll set the active tab
                clickedTab.setAttribute('aria-selected', true);
                clickedTab.setAttribute('tabindex', '0');
                // don't "focus".  this causes an unwanted blue border above and below the button which I can't figure out how to supress.
                //clickedTab.focus();

                // and then remove the hidden attribute from the corresponding active panel
                activePanel.removeAttribute('hidden');

                // TODO:
                // fire and event so the control surface gets updated
        }

        // make the default tab active:
        const activeTab = tabsList.querySelector('.nav-link.active');
        if (activeTab) {
                setActiveTab(activeTab);
        }

}
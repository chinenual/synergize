package com.chinenual.synergize;

import static com.chinenual.synergize.IntegrationTestSuite.driver;
import com.chinenual.synergize.pages.MainPage;
import com.chinenual.synergize.pages.PreferencesPage;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Order;
import org.junit.jupiter.api.Test;

import org.junit.jupiter.api.MethodOrderer.OrderAnnotation;
import org.junit.jupiter.api.TestMethodOrder;

@Order(2)
@TestMethodOrder(OrderAnnotation.class)
public class PreferencesPageTest {

    @Test
    @Order(1)
    public void prefs_Open() {
        MainPage.helpButton().click();
        MainPage.preferencesMenuItem().click();

        //String pageSource = driver.getPageSource();
        //System.out.println("PREFS PAGE SOURCE: " + pageSource);

        Assertions.assertEquals("Synergize Preferences", PreferencesPage.pageTitle(), "Prefs page title");
    }

}

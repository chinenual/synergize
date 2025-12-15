/*
 * Click nbfs://nbhost/SystemFileSystem/Templates/Licenses/license-default.txt to change this license
 * Click nbfs://nbhost/SystemFileSystem/Templates/Classes/Class.java to edit this template
 */
package com.chinenual.synergize.pages;

import static com.chinenual.synergize.IntegrationTestSuite.driver;
import io.appium.java_client.AppiumBy;
import org.openqa.selenium.NoSuchElementException;
import org.openqa.selenium.WebElement;

/**
 *
 * @author tynor
 */
public class MainPage {

    static String getChildText(WebElement el) {
        String result = "";
        for (WebElement child : el.findElements(AppiumBy.xpath("./child::*"))) {
            result += child.getText();
        }
        return result;
    }

    public static String pageTitle() {
        WebElement el = driver.findElement(AppiumBy.xpath("//XCUIElementTypeWindow"));
        return el.getText();
    }
             
    public static String synergyStatus() {
        try {
            WebElement el = driver.findElement(AppiumBy.accessibilityId("synergyName"));
            return getChildText(el);
        } catch (NoSuchElementException exc) {
            // expected when the text is empty - XCUI tree seems to omit the element of the text is empty
            return "";
        }
    }

    public static String controlSurfaceStatus() {
        try {
            WebElement el = driver.findElement(AppiumBy.accessibilityId("controlSurfaceName"));
            return getChildText(el);
        } catch (NoSuchElementException exc) {
            // expected when the text is empty - XCUI tree seems to omit the element of the text is empty
            return "";
        }
    }

    public static WebElement helpButton() {
        return driver.findElement(AppiumBy.accessibilityId("helpButton"));
    }
    
    public static WebElement preferencesMenuItem() {
        return driver.findElement(AppiumBy.accessibilityId("preferencesMenuItem"));
    }
}

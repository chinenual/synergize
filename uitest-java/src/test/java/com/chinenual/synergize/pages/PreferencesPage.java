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
public class PreferencesPage {


    public static String pageTitle() {
        WebElement el = driver.findElement(AppiumBy.xpath("//XCUIElementTypeWindow"));
        return el.getText();
    }
             
}

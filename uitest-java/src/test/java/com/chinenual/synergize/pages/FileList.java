/*
 * Click nbfs://nbhost/SystemFileSystem/Templates/Licenses/license-default.txt to change this license
 * Click nbfs://nbhost/SystemFileSystem/Templates/Classes/Class.java to edit this template
 */
package com.chinenual.synergize.pages;

import static com.chinenual.synergize.IntegrationTestSuite.driver;
import static com.chinenual.synergize.pages.BasePage.getChildText;
import io.appium.java_client.AppiumBy;
import org.openqa.selenium.NoSuchElementException;
import org.openqa.selenium.WebElement;

/**
 *
 * @author tynor
 */
public class FileList extends BasePage {

    public static WebElement fileLink(String name) {
        //WebElement list = driver.findElement(AppiumBy.accessibilityId("CRTfiles"));
        WebElement el = driver.findElement(AppiumBy.xpath("//XCUIElementTypeStaticText[@value=\"INTERNAL\"]"));
        //for (WebElement child : el.findElements(AppiumBy.xpath("./child::*"))) {
        //    System.out.println("child: " + child.toString());
        //    if (child.getAttribute("value") == name) {
        //        return child;
        //    }
        //}
        return el;
    }

    public static String crt_path() {
        try {
            WebElement el = driver.findElement(AppiumBy.accessibilityId("crt_path"));
            return getChildText(el);
        } catch (NoSuchElementException exc) {
            // expected when the text is empty - XCUI tree seems to omit the element of the text is empty
            return "";
        }
    }

}

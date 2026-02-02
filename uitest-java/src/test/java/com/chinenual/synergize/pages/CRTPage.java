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
public class CRTPage extends BasePage {

    public static WebElement editButton() {
        WebElement el = driver.findElement(AppiumBy.accessibilityId("editCRTButton"));
        return el;
    }
    
    public static String slotVoicename(int slot) {
        WebElement el = driver.findElement(AppiumBy.accessibilityId("crt_voicename_" + slot));
        return getChildText(el);
    }

    public static WebElement slotAddButton(int slot) {
        try {
            WebElement el = driver.findElement(AppiumBy.accessibilityId("crt_slot_add_button_" + slot));
            return el;
        } catch (NoSuchElementException ex) {
            return null;
        }
    }
    
    public static WebElement slotClearButton(int slot) {
        try {
            WebElement el = driver.findElement(AppiumBy.accessibilityId("crt_slot_clear_button_" + slot));
            return el;
        } catch (NoSuchElementException ex) {
            return null;
        }
    }

}

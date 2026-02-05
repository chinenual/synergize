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

    public static WebElement slotVoice(int slot) {
        try {
            // HACK: SPAN-CLICK-CHILD
            // Can't click directly on a <span> node with the accessibility label - need to navigate to the child node 
            // (which apparenty corresponds to the actual <span>): 
            //
            //   HTML:
            //    <span aria-label="crt_voicename1" onclick=xxxx ...></span>
            //
            //   XPath tree:      
            //    <XCUIElementTypeCell elementType="75" identifier="" value="" label="" title="" placeholderValue="" enabled="true" selected="false" x="1554" y="875" width="230" height="31">
            //        <XCUIElementTypeGroup elementType="3" identifier="" value="" label="crt_voicename_3" title="crt_voicename_3" placeholderValue="" enabled="true" selected="false" x="1559" y="882" width="62" height="16">
            //            <XCUIElementTypeStaticText elementType="48" identifier="" value="RRHODES" label="" title="" placeholderValue="" enabled="true" selected="false" x="1559" y="882" width="62" height="16"></XCUIElementTypeStaticText>
            // 
            WebElement el = driver.findElement(AppiumBy.xpath("//*[@label=\"crt_voicename_" + slot+"\"]/child::*"));
            return el;
        } catch (NoSuchElementException ex) {
            return null;
        }
    }

    public static String slotVoiceName(int slot) {
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

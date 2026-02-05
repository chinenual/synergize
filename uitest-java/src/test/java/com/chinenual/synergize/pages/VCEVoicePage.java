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
public class VCEVoicePage extends BasePage {


    public static String vce_crt_name() {
       
        WebElement el = driver.findElement(AppiumBy.accessibilityId("vce_crt_name"));  
        return getChildText(el);
    }
    
    public static WebElement backToCRT() {
       try {
           DUMP_PAGE();
           // HACK: see SPAN-CLICK-CHILD - except this time there is an extra <span> involved... so get child of the child?
           //    <span id="backToCRT">
           //      <span id="vce_crt_name">XXXX ...
           //
            WebElement el = driver.findElement(AppiumBy.xpath("//*[@label=\"backToCRT\"]/child::*/child::*"));
            return el;
        } catch (NoSuchElementException ex) {
            return null;
        }
    }

    public static String vce_name() {
       
        WebElement el = driver.findElement(AppiumBy.accessibilityId("vce_name"));  
        return getChildText(el);
    }

}

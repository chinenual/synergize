/*
 * Click nbfs://nbhost/SystemFileSystem/Templates/Licenses/license-default.txt to change this license
 * Click nbfs://nbhost/SystemFileSystem/Templates/Classes/Class.java to edit this template
 */
package com.chinenual.synergize.pages;

import com.chinenual.synergize.IntegrationTestSuite;
import io.appium.java_client.AppiumBy;
import org.openqa.selenium.WebElement;

/**
 *
 * @author tynor
 */
public class BasePage {

    static String getChildText(WebElement el) {
        String result = "";
        for (WebElement child : el.findElements(AppiumBy.xpath("./child::*"))) {
            result += child.getText();
        }
        return result;
    }

    static void DUMP_PAGE() {
        String pageSource = IntegrationTestSuite.driver.getPageSource();
        System.out.println("-----------------------------------\nPAGE SOURCE: " + pageSource + "\n-------------------------------");

    }

}

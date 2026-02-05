/*
 * Click nbfs://nbhost/SystemFileSystem/Templates/Licenses/license-default.txt to change this license
 * Click nbfs://nbhost/SystemFileSystem/Templates/Classes/Class.java to edit this template
 */
package com.chinenual.synergize;

import static com.chinenual.synergize.IntegrationTestSuite.driver;
import io.appium.java_client.AppiumBy;
import java.awt.image.BufferedImage;
import java.io.ByteArrayInputStream;
import java.io.File;
import java.text.SimpleDateFormat;
import java.util.Date;
import javax.imageio.ImageIO;
import org.junit.jupiter.api.extension.ExtensionContext;
import org.junit.jupiter.api.extension.TestWatcher;
import org.openqa.selenium.OutputType;
import org.openqa.selenium.WebElement;

/**
 *
 * @author tynor
 */
public class ScreenshotOnFailureExtension implements TestWatcher {

    @Override
    public void testAborted(ExtensionContext context, Throwable cause) {
        System.out.println("****** TEST ABORTED ****** " + context.getDisplayName() + (context.getExecutionException().isPresent() ? " true " : " false ") + cause.toString());
        captureWindowScreenshot(context);
    }

    @Override
    public void testFailed(ExtensionContext context, Throwable cause) {
        System.out.println("****** TEST FAILED ****** " + context.getDisplayName() + (context.getExecutionException().isPresent() ? " true " : " false ") + cause.toString());
        captureWindowScreenshot(context);
    }


    private void captureWindowScreenshot(ExtensionContext context) {
        String testClassName = context.getTestClass().orElseThrow().getSimpleName(); // e.g., "LoginTests"
        String testMethodName = context.getTestMethod().orElseThrow().getName(); // e.g., "testInvalidCredentials"
        String timestamp = new SimpleDateFormat("yyyyMMdd_HHmmss").format(new Date());
        String screenshotPath = String.format("screenshots/%s_%s_%s.png",
                timestamp, testClassName, testMethodName);

        System.out.println("****** CAPTURE SCREENSHOT ******");
        String pageSource = driver.getPageSource();
        System.out.println("PAGE SOURCE AT SCREENSHOT "+screenshotPath+": " + pageSource);



        WebElement el = IntegrationTestSuite.driver.findElement(AppiumBy.xpath("//XCUIElementTypeWindow"));

        try {
            byte[] img = el.getScreenshotAs(OutputType.BYTES);
            //byte[] img = IntegrationTestSuite.driver.getScreenshotAs(OutputType.BYTES);

            BufferedImage fullImg = ImageIO.read(new ByteArrayInputStream(img));
            ImageIO.write(fullImg, "png", new File(screenshotPath));

            System.out.println("*** Screenshot saved to: " + screenshotPath);
        } catch (Throwable ex) {
            System.err.println("Failed to capture screenshot: " + ex.getMessage());

        }

    }

}

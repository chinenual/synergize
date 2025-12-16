/*
 * Click nbfs://nbhost/SystemFileSystem/Templates/Licenses/license-default.txt to change this license
 * Click nbfs://nbhost/SystemFileSystem/Templates/Classes/Class.java to edit this template
 */
package com.chinenual.synergize;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.OpenOption;
import java.nio.file.Path;
import java.nio.file.StandardOpenOption;
import java.text.SimpleDateFormat;
import java.util.Date;
import org.junit.jupiter.api.extension.AfterTestExecutionCallback;
import org.junit.jupiter.api.extension.ExtensionContext;
import org.junit.jupiter.api.extension.TestWatcher;
import org.openqa.selenium.OutputType;
import org.openqa.selenium.TakesScreenshot;

/**
 *
 * @author tynor
 */
public class ScreenshotOnFailureExtension implements TestWatcher {

    @Override
    public void testAborted(ExtensionContext context, Throwable cause) {
        System.out.println("****** TEST ABORTED ****** " + context.getDisplayName() + (context.getExecutionException().isPresent() ? " true " : " false ") + cause.toString());
        captureScreenshot(context);
    }
    @Override
    public void testFailed(ExtensionContext context, Throwable cause) {
        System.out.println("****** TEST FAILED ****** " + context.getDisplayName() + (context.getExecutionException().isPresent() ? " true " : " false ") + cause.toString());
        captureScreenshot(context);
    }

    private void captureScreenshot(ExtensionContext context) {
        System.out.println("****** CAPTURE SCREENSHOT ******");
        String testClassName = context.getTestClass().orElseThrow().getSimpleName(); // e.g., "LoginTests"
        String testMethodName = context.getTestMethod().orElseThrow().getName(); // e.g., "testInvalidCredentials"
        String timestamp = new SimpleDateFormat("yyyyMMdd_HHmmss").format(new Date());
        String screenshotPath = String.format("screenshots/%s_%s_%s.png",
                timestamp, testClassName, testMethodName);

        try {
            byte[] img = IntegrationTestSuite.driver.getScreenshotAs(OutputType.BYTES);
            //File screenshotFile = ((TakesScreenshot) IntegrationTestSuite.driver).getScreenshotAs(OutputType.FILE);
            //File destFile = new File(screenshotPath);
            //copyFile(screenshotFile, destFile); // Save screenshot
            //System.out.println("Screenshot saved to: " + destFile.getAbsolutePath());
            Files.write(Path.of(screenshotPath), img, StandardOpenOption.CREATE, StandardOpenOption.WRITE);
            System.out.println("*** Screenshot saved to: " + screenshotPath);
        } catch (Throwable ex) {
            System.err.println("Failed to capture screenshot: " + ex.getMessage());

        }
    }

}

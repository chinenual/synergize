package com.chinenual.synergize;

import io.appium.java_client.AppiumBy;
import org.junit.jupiter.api.Assertions;
import io.appium.java_client.mac.Mac2Driver;
import io.appium.java_client.service.local.AppiumDriverLocalService;

import java.io.File;
import java.time.Duration;
import org.junit.jupiter.api.ClassOrderer;
import org.junit.jupiter.api.TestClassOrder;
import org.junit.platform.suite.api.AfterSuite;
import org.junit.platform.suite.api.BeforeSuite;
import org.junit.platform.suite.api.SelectPackages;
import org.junit.platform.suite.api.Suite;
import org.openqa.selenium.Dimension;
import org.openqa.selenium.Point;
import org.openqa.selenium.WebElement;
import org.openqa.selenium.remote.DesiredCapabilities;

/**
 * Integration (E2E) test for Synergize.
 */
@Suite
@SelectPackages("com.chinenual.synergize")
public class IntegrationTestSuite {

    public static Mac2Driver driver;

    @BeforeSuite
    public static void launchApp() {
        try {
            AppiumDriverLocalService service = AppiumDriverLocalService.buildDefaultService();
            service.start();
            System.out.println("Started Appium service: " + service.getUrl());

            DesiredCapabilities capabilities = new DesiredCapabilities();
            // App path must be absolute - appium looks for relative paths relative to the driver tmp directory
            String path = new File("../bin/Synergize.app").getCanonicalPath();
            // Appium will run the app in an appium driver specific tmp directory.
            // We need the app to run locally so it can find test files via relative paths:
            String[] args = {"-MOCKSYNIO", "-SERIALVERBOSE", "-CWD", System.getProperty("user.dir")};
            System.out.println("APP PATH: " + path);
            capabilities.setCapability("appium:appPath", path);
            capabilities.setCapability("appium:arguments", args);
            capabilities.setCapability("appium:bundleId", "com.chinenual.synergize");

            driver = new Mac2Driver(capabilities);

            String pageSource = driver.getPageSource();
            System.out.println("PAGE SOURCE: " + pageSource);

        } catch (Throwable exc) {
            System.err.println("ERROR: " + exc.toString());
            Assertions.fail(exc);
            if (driver != null) {
                
                System.err.println("INFO: driver.quit()");
                driver.quit();
            }
            System.exit(1);
        }
    }

    @AfterSuite
    public static void teardown() {
//        System.err.println("SLEEPING 20 min");
//        try {
//            Thread.sleep(Duration.ofMinutes(20));
//        } catch (InterruptedException ex) {
//            System.getLogger(IntegrationTestSuite.class.getName()).log(System.Logger.Level.ERROR, (String) null, ex);
//        }
        if (driver != null) {
            System.err.println("INFO: driver.quit()");
            driver.quit();
        }
    }

}

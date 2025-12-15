package com.chinenual.synergize;

import org.junit.jupiter.api.Assertions;
import io.appium.java_client.mac.Mac2Driver;
import io.appium.java_client.service.local.AppiumDriverLocalService;

import java.io.File;
import org.junit.platform.suite.api.AfterSuite;
import org.junit.platform.suite.api.BeforeSuite;
import org.junit.platform.suite.api.SelectPackages;
import org.junit.platform.suite.api.Suite;
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
            String path = new File("../bin/Synergize.dev.app").getCanonicalPath();
            String[] args = {"-MOCKSYNIO", "-SERIALVERBOSE"};
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
            System.exit(1);
        }
    }

    @AfterSuite
    public static void teardown() {
        if (driver != null) {
            driver.quit();
        }
    }

}

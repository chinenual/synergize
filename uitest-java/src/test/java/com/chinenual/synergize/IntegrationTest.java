package com.chinenual.synergize;

import com.chinenual.synergize.pages.MainPage;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.MethodOrderer.OrderAnnotation;
import org.junit.jupiter.api.Order;
import org.junit.jupiter.api.Test;

import io.appium.java_client.AppiumBy;
import io.appium.java_client.mac.Mac2Driver;
import io.appium.java_client.service.local.AppiumDriverLocalService;

import java.io.File;
import java.net.URI;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.TestMethodOrder;
import org.openqa.selenium.NoSuchElementException;
import org.openqa.selenium.WebElement;
import org.openqa.selenium.remote.DesiredCapabilities;

/**
 * Integration (E2E) test for Synergize.
 */
@TestMethodOrder(OrderAnnotation.class)
public class IntegrationTest {

    public static Mac2Driver driver;

    @BeforeAll
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

    @AfterAll
    public static void teardown() {
        if (driver != null) {
            driver.quit();
        }
    }

    @Test
    @Order(1)
    public void initialPageTitle() {
        Assertions.assertEquals("Synergize", MainPage.pageTitle());
    }
    
    @Test
    @Order(2)
    public void initialSynergyStatus() {
        Assertions.assertEquals("not connected", MainPage.synergyStatus());
    }

    @Test
    @Order(3)
    public void initialControlSurfaceStatus() {
        Assertions.assertEquals("", MainPage.controlSurfaceStatus());
    }
}

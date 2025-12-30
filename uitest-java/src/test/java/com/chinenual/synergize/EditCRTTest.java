package com.chinenual.synergize;

import static com.chinenual.synergize.IntegrationTestSuite.driver;
import com.chinenual.synergize.pages.FileList;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Order;
import org.junit.jupiter.api.Test;

import org.junit.jupiter.api.MethodOrderer.OrderAnnotation;
import org.junit.jupiter.api.TestInstance;
import org.junit.jupiter.api.TestInstance.Lifecycle;
import org.junit.jupiter.api.TestMethodOrder;
import org.junit.jupiter.api.extension.ExtendWith;
import org.openqa.selenium.WebElement;

@Order(4)
@ExtendWith(ScreenshotOnFailureExtension.class)
@TestMethodOrder(OrderAnnotation.class)
@TestInstance(Lifecycle.PER_CLASS) 
public class EditCRTTest {

    @BeforeAll
    public void init() {
        String pageSource = driver.getPageSource();
        System.out.println("PAGE SOURCE: " + pageSource);
    }

    @Test
    @Order(1)
    public void load_INTERNAL_CRT() {
        WebElement el = FileList.fileLink("INTERNAL");
        el.click();
        Assertions.assertEquals("INTERNAL", FileList.crt_path());
    }
    
    @Test
    @Order(2)
    public void forceFail() {
       Assertions.assertEquals(true,false);
    }

}

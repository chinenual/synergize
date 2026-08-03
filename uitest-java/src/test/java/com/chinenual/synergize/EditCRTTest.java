package com.chinenual.synergize;

import static com.chinenual.synergize.IntegrationTestSuite.driver;
import com.chinenual.synergize.pages.BasePage;
import com.chinenual.synergize.pages.CRTPage;
import com.chinenual.synergize.pages.FileList;
import com.chinenual.synergize.pages.VCEVoicePage;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Order;
import org.junit.jupiter.api.Test;

import org.junit.jupiter.api.MethodOrderer.OrderAnnotation;
import org.junit.jupiter.api.TestInstance;
import org.junit.jupiter.api.TestInstance.Lifecycle;
import org.junit.jupiter.api.TestMethodOrder;
import org.junit.jupiter.api.extension.ExtendWith;
import org.openqa.selenium.By;
import org.openqa.selenium.Keys;
import org.openqa.selenium.WebElement;
import org.openqa.selenium.interactions.Actions;

@Order(4)
@ExtendWith(ScreenshotOnFailureExtension.class)
@TestMethodOrder(OrderAnnotation.class)
@TestInstance(Lifecycle.PER_CLASS) 
public class EditCRTTest {

    //@BeforeAll
    //public void init() {
    //    String pageSource = driver.getPageSource();
    //    System.out.println("PAGE SOURCE: " + pageSource);
    //}

    @Test
    @Order(1)
    public void load_INTERNAL_CRT() {
        WebElement el = FileList.fileLink("INTERNAL");
        el.click();
        Assertions.assertEquals("INTERNAL", FileList.crt_path());
        
        el = CRTPage.slotClearButton(1);
        Assertions.assertNull(el, "clear buttons should be hidden");
        el = CRTPage.slotAddButton(1);
        Assertions.assertNull(el, "add buttons should be hidden");
    }
    
    @Test
    @Order(2)
    public void checkVoiceSlots() {
       // sanity check a few slots:
       Assertions.assertEquals("G7S", CRTPage.slotVoiceName(1));
       Assertions.assertEquals("HORNSXX", CRTPage.slotVoiceName(2));
       Assertions.assertEquals("RRHODES", CRTPage.slotVoiceName(3));
    }
    
    @Test
    @Order(3)
    public void viewVoice() {       
       WebElement el = CRTPage.slotVoice(3);
       el.click();
       Assertions.assertEquals("INTERNAL", VCEVoicePage.vce_crt_name());
       Assertions.assertEquals("RRHODES",  VCEVoicePage.vce_name());               
    }
    
    @Test
    @Order(4)
    public void reloadINTERNAL_CRT() {      
       VCEVoicePage.backToCRT().click();
               
       // sanity check a few slots:
       Assertions.assertEquals("G7S", CRTPage.slotVoiceName(1));
       Assertions.assertEquals("HORNSXX", CRTPage.slotVoiceName(2));
       Assertions.assertEquals("RRHODES", CRTPage.slotVoiceName(3));
    }
    
    @Test
    @Order(5)
    public void enableEdit() {
        WebElement el = CRTPage.editButton();
        el.click();
        
        el = CRTPage.slotClearButton(1);
        Assertions.assertNotNull(el, "clear buttons should be visible");
        el = CRTPage.slotAddButton(1);
        Assertions.assertNotNull(el, "add buttons should be visible");

       // sanity check a few slots:
       Assertions.assertEquals("G7S", CRTPage.slotVoiceName(1));
       Assertions.assertEquals("HORNSXX", CRTPage.slotVoiceName(2));
       Assertions.assertEquals("RRHODES", CRTPage.slotVoiceName(3));
    }
  
    @Test
    @Order(6)
    public void editViewVoice() {       
       WebElement el = CRTPage.slotVoice(3);
       el.click();
       Assertions.assertEquals("INTERNAL", VCEVoicePage.vce_crt_name());
       Assertions.assertEquals("RRHODES",  VCEVoicePage.vce_name());               
    }
    
    @Test
    @Order(7)
    public void editReloadINTERNAL_CRT() {      
       VCEVoicePage.backToCRT().click();
               
       // sanity check a few slots:
       Assertions.assertEquals("G7S", CRTPage.slotVoiceName(1));
       Assertions.assertEquals("HORNSXX", CRTPage.slotVoiceName(2));
       Assertions.assertEquals("RRHODES", CRTPage.slotVoiceName(3));
    }
    
    @Test
    @Order(8)
    public void deleteVoice() {
        WebElement el = CRTPage.slotClearButton(24);
        el.click();
        Assertions.assertEquals("", CRTPage.slotVoiceName(24), "slot name is cleared");
        el = CRTPage.slotClearButton(24);
        Assertions.assertNotNull(el, "clear buttons (24) should be visible");
        el = CRTPage.slotAddButton(24);
        Assertions.assertNotNull(el, "add buttons (24) should be visible");
    }
    
    @Test
    @Order(9)
    public void addVoice() {
        WebElement el = CRTPage.slotAddButton(23);
        el.click(); 
                
        el = BasePage.fileDialog();
        // cmd-shift-g to allow us to send an absolute path:
        Actions a = new Actions(IntegrationTestSuite.driver);
        a.keyDown(Keys.COMMAND)
                .keyDown(Keys.SHIFT)
                .sendKeys("g")
                .build()
                .perform();
        
        BasePage.DUMP_PAGE();
        Assertions.assertEquals("COWBELL", CRTPage.slotVoiceName(23), "slot name is new voice");
    }

}

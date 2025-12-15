package com.chinenual.synergize;

import com.chinenual.synergize.pages.MainPage;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Order;
import org.junit.jupiter.api.Test;

import org.junit.jupiter.api.MethodOrderer.OrderAnnotation;
import org.junit.jupiter.api.TestMethodOrder;

@DisplayName("MainPageTest display name")
@Order(1)
@TestMethodOrder(OrderAnnotation.class)
public class MainPageTest {

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

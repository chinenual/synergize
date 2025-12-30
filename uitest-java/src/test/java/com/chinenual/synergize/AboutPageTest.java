package com.chinenual.synergize;

import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.DisplayNameGenerator;
import org.junit.jupiter.api.IndicativeSentencesGeneration;
import org.junit.jupiter.api.Order;
import org.junit.jupiter.api.Test;

import org.junit.jupiter.api.MethodOrderer.OrderAnnotation;
import org.junit.jupiter.api.TestMethodOrder;
import org.junit.jupiter.api.extension.ExtendWith;

@Order(3)
@ExtendWith(ScreenshotOnFailureExtension.class)
@TestMethodOrder(OrderAnnotation.class)
public class AboutPageTest {

    @Test
    @Order(1)
    public void abouttest1() {
        //Assertions.assertEquals(true,false);
    }
   
}

package com.balance.backend.controller;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class Status {

    @GetMapping("/status")
    public String status() {
        return "OK";
    }

}

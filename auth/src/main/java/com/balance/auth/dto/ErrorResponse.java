package com.balance.auth.dto;

import lombok.Data;

import java.sql.Timestamp;

@Data
public class ErrorResponse {

    private String message;
    private Timestamp timestamp = new Timestamp(System.currentTimeMillis());
    private String correlationId;

}

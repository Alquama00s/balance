package com.balance.auth.exceptions;

public class DuplicateRecordException extends RuntimeException {
    public DuplicateRecordException(String recordType, String fieldName, String fieldValue) {
        super(String.format("%s with %s '%s' already exists", recordType, fieldName, fieldValue));
    }
}

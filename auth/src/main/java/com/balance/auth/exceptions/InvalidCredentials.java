package com.balance.auth.exceptions;

public class InvalidCredentials extends RuntimeException {
    public InvalidCredentials() {
        super("Invalid username or password");
    }
}

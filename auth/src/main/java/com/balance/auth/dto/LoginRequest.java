package com.balance.auth.dto;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import lombok.Data;

@Data
public class LoginRequest {
    @NotBlank(message = "Username must be valid")
    String username;
    @NotNull
    @Size(min = 6,message = "invalid password")
    String password;
}

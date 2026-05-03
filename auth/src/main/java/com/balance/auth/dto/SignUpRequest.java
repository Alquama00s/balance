package com.balance.auth.dto;

import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;
import lombok.Data;

@Data
public class SignUpRequest {
    @NotBlank(message = "Username must be valid")
    private String firstName;
    private String lastName;
    @Email
    private String email;
    @Size(min = 6,message = "Password must be at least 6 characters long")
    private String password;
    private String phoneNumber;
}

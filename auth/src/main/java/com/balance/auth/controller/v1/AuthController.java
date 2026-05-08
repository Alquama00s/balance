package com.balance.auth.controller.v1;

import com.balance.auth.dto.LoginRequest;
import com.balance.auth.dto.LoginResponse;
import com.balance.auth.dto.SignUpRequest;
import com.balance.auth.dto.SignUpResponse;
import com.balance.auth.service.AuthService;
import com.nimbusds.jose.JOSEException;
import jakarta.validation.Valid;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/v1/auth")
public class AuthController {

    private AuthService authService;

    @Autowired
    public void setAuthService(AuthService authService) {
        this.authService = authService;
    }

    @PostMapping("/signup")
    public SignUpResponse signUp(@RequestBody @Valid SignUpRequest signUpRequest) {
        return authService.signUp(signUpRequest);
    }

    @PostMapping("/login")
    public LoginResponse login(@RequestBody @Valid LoginRequest loginRequest) throws JOSEException {
        return authService.login(loginRequest);
    }


}

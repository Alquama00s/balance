package com.balance.auth.controller.v1;

import com.balance.auth.dto.SignUpRequest;
import com.balance.auth.dto.SignUpResponse;
import com.balance.auth.service.AuthService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;

@Controller
@RequestMapping("/api/v1/auth")
public class AuthController {

    private AuthService authService;

    @Autowired
    public void setAuthService(AuthService authService) {
        this.authService = authService;
    }

    @PostMapping("/signup")
    public SignUpResponse signUp(SignUpRequest signUpRequest) {
        return authService.signUp(signUpRequest);
    }

}

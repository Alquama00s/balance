package com.balance.auth.service;

import com.balance.auth.dto.LoginRequest;
import com.balance.auth.dto.LoginResponse;
import com.balance.auth.dto.SignUpRequest;
import com.balance.auth.dto.SignUpResponse;
import com.balance.auth.entity.User;
import com.balance.auth.entity.UserProfile;
import com.balance.auth.exceptions.DuplicateRecordException;
import com.balance.auth.exceptions.InvalidCredentials;
import com.balance.auth.repository.UserProfileRepository;
import com.balance.auth.repository.UserRepository;
import com.nimbusds.jose.JOSEException;
import jakarta.transaction.Transactional;
import org.mindrot.jbcrypt.BCrypt;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class AuthService {

    private UserService userService;
    private UserProfileRepository userProfileRepository;



    private JWTService jwtService;

    @Autowired
    public void setUserRepository(UserService userService) {
        this.userService = userService;
    }
    @Autowired
    public void setUserProfileRepository(UserProfileRepository userProfileRepository) {
        this.userProfileRepository = userProfileRepository;
    }
    @Autowired
    public void setJwtService(JWTService jwtService) {
        this.jwtService = jwtService;
    }


    @Transactional
    public SignUpResponse signUp(SignUpRequest signUpRequest) {
        var existingUser = userService.findByUsername(signUpRequest.getEmail());
        if(existingUser.isPresent()){
            throw new DuplicateRecordException("User","email", signUpRequest.getEmail());
        }
        var user = new User();
        user.setEmail(signUpRequest.getEmail());
        user.setUsername(signUpRequest.getEmail());
        user.setPasswordHash(hashPassword(signUpRequest.getPassword()));
        user.setCreatedBy(userService.getAdminUser().getId());
        user = userService.save(user);

        UserProfile userProfile = new UserProfile();
        userProfile.setFirstName(signUpRequest.getFirstName());
        userProfile.setLastName(signUpRequest.getLastName());
        userProfile.setPhoneE164(signUpRequest.getPhoneNumber());
        userProfile.setUser(user);

        userProfileRepository.save(userProfile);

        var signUpRes = new SignUpResponse();
        signUpRes.setUserName(user.getUsername());
        signUpRes.setEmail(user.getEmail());
        return signUpRes;
    }

    public LoginResponse login(LoginRequest loginRequest) throws JOSEException {
        var user = userService.findByUsername(loginRequest.getUsername()).orElseThrow();
        if (checkPassword(loginRequest.getPassword(), user.getPasswordHash())) {
            var loginRes = new LoginResponse();
            loginRes.setToken(jwtService.generateToken(user));
            return loginRes;
        } else {
            throw new InvalidCredentials();
        }
    }


    private String hashPassword(String password) {
        return BCrypt.hashpw(password, BCrypt.gensalt());
    }

    private boolean checkPassword(String password, String hash) {
        return BCrypt.checkpw(password, hash);
    }

}

package com.balance.auth.service;

import com.balance.auth.dto.SignUpRequest;
import com.balance.auth.dto.SignUpResponse;
import com.balance.auth.entity.User;
import com.balance.auth.entity.UserProfile;
import com.balance.auth.repository.UserProfileRepository;
import com.balance.auth.repository.UserRepository;
import jakarta.transaction.Transactional;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class AuthService {

    private UserRepository userRepository;
    private UserProfileRepository userProfileRepository;
    @Autowired
    public void setUserRepository(UserRepository userRepository) {
        this.userRepository = userRepository;
    }
    @Autowired
    public void setUserProfileRepository(UserProfileRepository userProfileRepository) {
        this.userProfileRepository = userProfileRepository;
    }

    @Transactional
    public SignUpResponse signUp(SignUpRequest signUpRequest) {

        var user = new User();
        user.setEmail(signUpRequest.getEmail());
        user.setUsername(signUpRequest.getEmail());
        user.setPasswordHash(signUpRequest.getPassword());
        user = userRepository.save(user);

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




}

package com.balance.auth.service;

import com.balance.auth.entity.User;
import com.balance.auth.repository.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class UserService {

    @Autowired
    private UserRepository userRepository;


    public User save(User user) {
        return userRepository.save(user);
    }

    public User getAdminUser() {
        var admin = userRepository.findByUsername("admin");
        if(admin.isEmpty()){
            var u = new User();
            u.setEmail("admin@admin.com");
            u.setUsername("admin");
            u.setPasswordHash("jughuyghu");
            return save(u);
        }
        return admin.get();
    }

}


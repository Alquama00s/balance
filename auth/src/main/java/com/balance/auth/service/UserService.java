package com.balance.auth.service;

import com.balance.auth.entity.User;
import com.balance.auth.repository.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Optional;

@Service
public class UserService {

    @Autowired
    private UserRepository userRepository;

    private User adminUser;

    public User save(User user) {
        return userRepository.save(user);
    }

    public Optional<User> findByUsername(String username) {
        return userRepository.findByUsername(username);
    }

    public User getAdminUser() {
        if(adminUser != null){
            return adminUser;
        }
        var admin = userRepository.findByUsername("admin");
        if(admin.isEmpty()){
            var u = new User();
            u.setEmail("admin@admin.com");
            u.setUsername("admin");
            u.setPasswordHash("jughuyghu");
            return save(u);
        }
        return adminUser = admin.get();
    }

}


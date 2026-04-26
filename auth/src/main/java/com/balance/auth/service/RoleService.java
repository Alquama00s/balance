package com.balance.auth.service;
import com.balance.auth.entity.Privilege;
import com.balance.auth.entity.Role;
import com.balance.auth.entity.User;
import com.balance.auth.repository.RoleRepository;
import jakarta.annotation.PostConstruct;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.HashSet;
import java.util.Set;
import java.util.stream.Collectors;

@Service
public class RoleService {

    @Autowired
    RoleRepository roleRepository;

    @Autowired
    UserService userService;

    private Set<Privilege> getDefaultPrivileges(){
        final Set<String> actions = new HashSet<>();
        actions.add("create");
        actions.add("read");
        actions.add("update");
        actions.add("delete");
        final Set<String> entities = new HashSet<>();
        entities.add("transaction");
        entities.add("accounts");

        return entities.stream()
                .flatMap(e-> actions.stream().map(a->String.join("_",e,a)))
                .map(this::fromDefaultPrivilege)
                .collect(Collectors.toSet());
    }

    private Privilege fromDefaultPrivilege(String privilege){
        var admin = userService.getAdminUser();
        var p = new Privilege();
        p.setName(privilege);
        p.setDescription("auto generated");
        p.setCreatedBy(admin.getId());
        return p;
    }

    @PostConstruct
    private void initialise(){
        var role = roleRepository.findByName("user");
        var admin = userService.getAdminUser();
        if(role.isEmpty()){
            var r = new Role();
            r.setPrivileges(getDefaultPrivileges());
            r.setName("user");
            r.setDescription("auto generated");
            r.setCreatedBy(admin.getId());
            roleRepository.save(r);
        }
    }

}

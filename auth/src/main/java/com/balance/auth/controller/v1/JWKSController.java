package com.balance.auth.controller.v1;


import com.nimbusds.jose.KeySourceException;
import com.nimbusds.jose.jwk.JWKMatcher;
import com.nimbusds.jose.jwk.JWKSelector;
import com.nimbusds.jose.jwk.source.JWKSource;
import com.nimbusds.jose.proc.SecurityContext;
import jakarta.annotation.PostConstruct;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.HashMap;
import java.util.Map;

@RestController
@RequestMapping("/oauth")
public class JWKSController {

    @Autowired
    private JWKSource<SecurityContext> jwkSource;


    HashMap<String, Object> publicKeys;


    @PostConstruct
    public void init() throws KeySourceException {
        publicKeys = new HashMap<>();
        var matcher = new JWKMatcher.Builder()
                .build();
        var keys = jwkSource.get(new JWKSelector(matcher), null);
        for (var i : keys) {
            publicKeys.put(i.getKeyID(), i.toPublicJWK().toJSONObject());
        }
    }

    @GetMapping("/jwks")
    public Map<String, Object> getAvailableKeys() throws KeySourceException {
        return Map.of("keys", publicKeys.values());
    }

}

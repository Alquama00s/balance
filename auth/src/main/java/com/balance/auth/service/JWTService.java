package com.balance.auth.service;

import com.balance.auth.entity.User;
import com.balance.auth.utills.JWTUtill;
import com.nimbusds.jose.JOSEException;
import com.nimbusds.jose.JWSAlgorithm;
import com.nimbusds.jose.JWSHeader;
import com.nimbusds.jose.KeySourceException;
import com.nimbusds.jose.crypto.RSASSASigner;
import com.nimbusds.jose.jwk.JWK;
import com.nimbusds.jose.jwk.JWKMatcher;
import com.nimbusds.jose.jwk.JWKSelector;
import com.nimbusds.jose.jwk.RSAKey;
import com.nimbusds.jose.jwk.source.JWKSource;
import com.nimbusds.jose.proc.SecurityContext;
import com.nimbusds.jwt.JWTClaimsSet;
import com.nimbusds.jwt.JWTParser;
import com.nimbusds.jwt.SignedJWT;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.security.SecureRandom;
import java.util.Random;

@Service
public class JWTService {
    private JWKSource<SecurityContext> jwkSource;
    private final JWKSelector allJwkSelector = new JWKSelector(new JWKMatcher.Builder().build());
    private final Random random = new SecureRandom();

    @Autowired
    public void setJwkSource(JWKSource<SecurityContext> jwkSource) {
        this.jwkSource = jwkSource;
    }


    public String generateToken(User user) throws JOSEException {
        var jwk = getRandomJWK();
        var signer = new RSASSASigner(jwk);
        JWTClaimsSet claimsSet = new JWTClaimsSet.Builder()
                .subject(user.getUsername())
                .claim("userId", user.getId())
                .claim("roles", user.getRoles())
                .build();
        JWSHeader header = new JWSHeader.Builder(JWSAlgorithm.RS256)
                .keyID(jwk.getKeyID())
                .build();
        SignedJWT signedJWT = new SignedJWT(header, claimsSet);
        signedJWT.sign(signer);
        return signedJWT.serialize();
    }

    private RSAKey getRandomJWK() throws KeySourceException {
        var keyList = jwkSource.get(allJwkSelector,null);
        return (RSAKey) keyList.get(random.nextInt(keyList.size()));
    }

}

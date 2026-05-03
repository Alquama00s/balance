package com.balance.auth.entity;


import jakarta.persistence.*;
import lombok.Data;
import lombok.EqualsAndHashCode;

import java.util.UUID;

@Entity
@Table(name = "user_profile")
@Data
@EqualsAndHashCode(callSuper = false)
public class UserProfile extends BaseAuditEntity {

    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    private UUID id;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "user_id", nullable = false, unique = true)
    private User user;

    @Column(name = "first_name", length = 100)
    private String firstName;

    @Column(name = "last_name", length = 100)
    private String lastName;

    @Column(name = "phone_e164", length = 20)
    private String phoneE164;

    @Column(name = "phone_verified")
    private Boolean phoneVerified = false;

    @Column(name = "avatar_url", columnDefinition = "TEXT")
    private String avatarUrl;
}


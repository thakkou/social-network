'use client';

import React, { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import { useSession } from "next-auth/react";

import ConfirmModal from "~/app/_components/ConfirmModal";
import ProfilePosts from "~/app/_components/ProfilePosts";
import Followers from "~/app/_components/Followers";
import { getProfileData } from "~/app/_services/crud/getProfile";
import PrivateProfile from "~/app/_components/PrivateProfile";
import { toggleFollow } from "~/app/_services/crud/follow";
import { useChat } from "~/app/_providers/chatProvider";
type TabItem = {
    label: string;
    Component: React.ComponentType<any>;
    props: any;
}

export default function Profile() {
    const params = useParams();
    const router = useRouter();
    const { data: session } = useSession();
    const userId = params?.id as string;

    // Redirect to /profile if viewing own profile
    useEffect(() => {
      if (userId && session?.user?.id && userId === session.user.id) {
        router.replace("/profile");
      }
    }, [userId, session?.user?.id, router]);

    const [activeTab, setActiveTab] = useState(0);
    const [isPrivate, setIsPrivate] = useState<boolean>(false);
    const [isLoading, setIsLoading] = useState(true);
    const [profileExists, setProfileExists] = useState(true);
    const [profile, setProfile] = useState<any>(null);
    const [isFollowLoading, setIsFollowLoading] = useState(false);
    const [showUnfollowConfirm, setShowUnfollowConfirm] = useState(false);
    const { selectChat } = useChat();

    // `tabs` depends on `profile`, so it must be declared after the state above
    const tabs: TabItem[] = [
        {
            label: "posts",
            Component: ProfilePosts,
            props: {
                posts: profile?.posts ?? []
            }
        },
        {
            label: "followers",
            Component: Followers,
            props: {
                users: profile?.followers ?? []
            }
        },
        {
            label: "following",
            Component: Followers,
            props: {
                users: profile?.following ?? []
            }
        }
    ];

    const CurrentTab = tabs[activeTab];

   useEffect(() => {
    if (!userId) return;

    const getData = async () => {
        setIsLoading(true);

        try {
            const res = await getProfileData(userId);
            if (res.success) {
                console.log(res.data)
                setProfile(res.data); // adjust to your API
                res.data.is_private==0 ? setIsPrivate(false) : setIsPrivate(true)
                setProfileExists(true);
            } else {
                setProfileExists(false);
            }
        } catch (err) {
            console.error(err);
            setProfileExists(false);
        } finally {
            setIsLoading(false);
        }
    };

    void getData();
}, [userId]);

const isPrivateBlocked =
    isPrivate &&
    profile?.following_status !== "accepted"

    // following_status can be "none", "pending", or "accepted"
    const followLabel =
        profile?.following_status === "accepted"
            ? "unfollow"
            : profile?.following_status === "pending"
            ? "requested"
            : "follow";

    const followIcon =
        profile?.following_status === "accepted"
            ? "ti-user-check"
            : profile?.following_status === "pending"
            ? "ti-clock"
            : "ti-user-plus";

    const handleFollowClick = async () => {
        if (!userId || isFollowLoading) return;

        // Show confirmation before unfollowing
        if (profile?.following_status === "accepted") {
            setShowUnfollowConfirm(true);
            return;
        }

        await executeFollow();
    };

    const executeFollow = async () => {
        if (!userId || isFollowLoading) return;
        setIsFollowLoading(true);

        try {
            const res = await toggleFollow(userId, profile.following_status);
            if ("error" in res) {
                console.error(res.error);
                return;
            }

            setProfile((prev: any) => ({
                ...prev,
                following_status: res.status,
            }));
        } catch (err) {
            console.error(err);
        } finally {
            setIsFollowLoading(false);
            setShowUnfollowConfirm(false);
        }
    };

    const confirmUnfollow = () => {
        void executeFollow();
    };

    const cancelUnfollow = () => {
        setShowUnfollowConfirm(false);
    };


    if (isLoading) {
        return (
            <main className="main">
                <div className="card" style={{ textAlign: 'center', padding: '32px' }}>
                    <p style={{ fontSize: '12px', color: 'var(--color-text-secondary)' }}>Loading profile...</p>
                </div>
            </main>
        );
    }


    if (!profileExists) {
        return (
            <main className="main">
                <div className="card" style={{ textAlign: 'center', padding: '32px' }}>
                    <p style={{ fontSize: '14px', fontWeight: 500, color: 'var(--color-text-primary)' }}>
                        This profile doesn't exist
                    </p>
                    <p style={{ fontSize: '11px', color: 'var(--color-text-secondary)', marginTop: '4px' }}>
                        The user you're looking for may have been removed or the link is incorrect.
                    </p>
                </div>
            </main>
        );
    }
if (isPrivateBlocked) {
    return (
        <main className="main">
            <PrivateProfile
            userId={userId}
                profile={{
                    firstname: profile.firstname,
                    lastname: profile.lastname,
                    nickname: profile.nickname,
                     following_status:profile.following_status,   
                    avatar: profile.avatar,
                }}
                onSendRequest={() => {
                    void handleFollowClick();
                }}
            />
        </main>
    );
}


    // Determine display name and initials
    const displayName = profile?.nickname || `${profile?.firstname || ''} ${profile?.lastname || ''}`.trim() || 'User';
    const avatarInitials = (profile?.firstname?.[0]?.toUpperCase() || '') + (profile?.lastname?.[0]?.toUpperCase() || '') || '?';
    const isFollowing = profile?.following_status === "accepted";

    return (<>
            <ConfirmModal
                open={showUnfollowConfirm}
                title="Unfollow user"
                message={`Are you sure you want to unfollow ${displayName}?`}
                confirmLabel="unfollow"
                confirmClass="btn-red"
                onConfirm={confirmUnfollow}
                onCancel={cancelUnfollow}
            />

            <main className="main">
            <div className="card">
                <div style={{ display: 'flex', alignItems:'flex-start', gap:'12px', marginBottom:'12px' }}>
                <div
                  className="av"
                  style={{
                    width: '52px',
                    height: '52px',
                    background: profile?.avatar ? `url(${profile.avatar}) center/cover` : '#EEEDFE',
                    color: '#534AB7',
                    fontSize: '16px',
                    overflow: 'hidden',
                  }}
                >
                  {!profile?.avatar && avatarInitials}
                </div>
                <div style={{ flex:1 }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap:'8px', flexWrap: 'wrap', marginBottom:'4px' }}>
                        <p style={{ fontSize:'14px', fontWeight:500, color:'var(--color-text-primary)' }}>{displayName}</p>
                        {profile?.nickname && (
                          <span className="tag tag-teal">@{profile.nickname}</span>
                        )}
                        <span className={`tag ${isPrivate ? 'tag-gray' : 'tag-purple'}`} id="profile-visibility-tag">{isPrivate ? 'private' : 'public'}</span>
                        <button
                            className="btn btn-g"
                            style={{ fontSize: '10px', display: 'flex', alignItems: 'center', gap: '3px' }}
                            id="follow-btn"
                            disabled={profile?.following_status === 'pending' || isFollowLoading}
                            onClick={() => void handleFollowClick()}
                        >
                            <i className={`ti ${followIcon}`} style={{ fontSize: '12px' }} aria-hidden="true"></i> {isFollowLoading ? '...' : followLabel}
                        </button>
                        <button
                            className="btn btn-p"
                            style={{ fontSize: '10px', display: 'flex', alignItems: 'center', gap: '3px' }}
                            onClick={() => {
                                selectChat({
                                    id: String(userId),
                                    type: "user",
                                    data: {
                                        other_user_id: Number(userId),
                                        display_name: displayName,
                                        avatar: profile?.avatar || "",
                                    },
                                });
                                router.push('/messages');
                            }}
                        >
                            <i className="ti ti-message" style={{ fontSize: '12px' }} aria-hidden="true"></i> message
                        </button>
                    </div>
                    {/* Show email + birth only if following */}
                    {isFollowing && (
                      <p style={{ fontSize:'11px', color:'var(--color-text-secondary)', marginBottom:'4px' }}>
                        {profile?.birthdate ? `Born ${profile.birthdate}` : ''}{profile?.birthdate && profile?.email ? ' · ' : ''}{profile?.email || ''}
                      </p>
                    )}
                    {profile?.aboutme && (
                      <p style={{ fontSize:'12px', color:'var(--color-text-primary)', lineHeight:1.5 }}>{profile.aboutme}</p>
                    )}
                </div>
                </div>
                <div className="divider"></div>
                <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3,1fr)', gap:'8px', textAlign: 'center' }}>
                    <div style={{ background:'var(--color-background-secondary)', padding:'8px' }}>
                        <p style={{ fontSize:'18px', fontWeight:500, color:'#D4537E' }}>{profile?.posts?.length || "0"}</p>
                        <p style={{ fontSize:'10px', color:'var(--color-text-tertiary)' }}>posts</p>
                    </div>
                    <div style={{ background:'var(--color-background-secondary)', padding:'8px', cursor: 'pointer' }} onClick={ () => setActiveTab(tabs.findIndex(o => o.label === 'followers')) }>
                        <p style={{ fontSize:'18px', fontWeight:500, color:'#534AB7' }}>{profile?.followers?.length || "0"}</p>
                        <p style={{ fontSize:'10px', color:'var(--color-text-tertiary)' }}>followers</p>
                    </div>
                    <div style={{ background: 'var(--color-background-secondary)', padding: '8px', cursor: 'pointer' }} onClick={() => setActiveTab(tabs.findIndex(o => o.label === 'following')) }>
                        <p style={{ fontSize:'18px', fontWeight:500, color:'#0F6E56' }}>{profile?.following?.length || "0"}</p>
                        <p style={{ fontSize:'10px', color:'var(--color-text-tertiary)' }}>following</p>
                    </div>
                </div>
            </div>

            <div style={{ display: 'flex', gap:0, border:'0.5px solid var(--color-border-tertiary)', background:'var(--color-background-primary)' }}>
                {tabs.map((tab, idx) => (
                    <div
                        key={idx}
                        className={`profile-tab ${activeTab === idx ? 'active-tab' : ''}`}
                        style={{
                            padding: '8px 16px',
                            fontSize: '11px',
                            cursor: 'pointer',
                            borderRight: idx + 1 < tabs.length ? '0.5px solid var(--color-border-tertiary)' : '',
                            color: activeTab === idx ? '#D4537E' : 'var(--color-text-secondary)',
                            borderBottom: activeTab === idx ? '2px solid #D4537E' : 'none' }}
                        onClick={() => setActiveTab(idx)}
                    >
                        {tab.label}
                    </div>
                ))}
            </div>

            {CurrentTab && <CurrentTab.Component {...CurrentTab.props} />}
        </main>
    </>);
}
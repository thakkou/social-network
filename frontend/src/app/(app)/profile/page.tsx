'use client';

import React, { useState } from "react";

import ProfilePosts from "../../_components/ProfilePosts";
import Followers from "../../_components/Followers";
import Following from "../../_components/Following";

interface TabItem {
    label: string;
    // This tells TS it is a valid React component
    Component: React.ComponentType<any>;
}

export default function Profile() {
    const tabs: TabItem[] = [
        { label: 'posts', Component: ProfilePosts },
        { label: 'followers', Component: Followers },
        { label: 'following', Component: Following },
    ];

    const [activeTab, setActiveTab] = useState(0);
    const [isPrivate, setIsPrivate] = useState(false); // what is a private profile ?!

    const CurrentTabComponent = tabs[activeTab]?.Component as React.ComponentType; // the one in TabItem didn't work !

    return (
        <main className="main">
            <div className="card">
                <div style={{ display: 'flex', alignItems:'flex-start', gap:'12px', marginBottom:'12px' }}>
                <div className="av" style={{ width:'52px', height:'52px', background:'#EEEDFE', color:'#534AB7', fontSize:'16px' }}>AK</div>
                <div style={{ flex:1 }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap:'8px', flexWrap: 'wrap', marginBottom:'4px' }}>
                        <p style={{ fontSize:'14px', fontWeight:500, color:'var(--color-text-primary)' }}>Amir Kader</p>
                        <span className="tag tag-teal">@amirkader</span>
                        <span className={`tag ${isPrivate ? 'tag-gray' : 'tag-purple'}`} id="profile-visibility-tag">{isPrivate ? 'private' : 'public'}</span>
                        <button className="btn btn-g" style={{ fontSize: '10px', display: 'flex', alignItems: 'center', gap: '3px', marginLeft: 'auto' }} id="visibility-btn" onClick={() => setIsPrivate(!isPrivate)}>
                            <i className={`ti ${isPrivate ? 'ti-lock-open' : 'ti-lock'}`} style={{ fontSize: '12px' }} aria-hidden="true"></i> make {isPrivate ? 'public' : 'private'}
                        </button>
                    </div>
                    <p style={{ fontSize:'11px', color:'var(--color-text-secondary)', marginBottom:'4px' }}>Born 1998-07-14 · amir@example.com</p>
                    <p style={{ fontSize:'12px', color:'var(--color-text-primary)', lineHeight:1.5 }}>Full-stack dev working on distributed systems. Passionate about Go, open source, and clean architecture.</p>
                </div>
                </div>
                <div className="divider"></div>
                <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3,1fr)', gap:'8px', textAlign: 'center' }}>
                    <div style={{ background:'var(--color-background-secondary)', padding:'8px' }}>
                        <p style={{ fontSize:'18px', fontWeight:500, color:'#D4537E' }}>12</p>
                        <p style={{ fontSize:'10px', color:'var(--color-text-tertiary)' }}>posts</p>
                    </div>
                    <div style={{ background:'var(--color-background-secondary)', padding:'8px', cursor: 'pointer' }} onClick={ () => setActiveTab(tabs.findIndex(o => o.label === 'followers')) }>
                        <p style={{ fontSize:'18px', fontWeight:500, color:'#534AB7' }}>48</p>
                        <p style={{ fontSize:'10px', color:'var(--color-text-tertiary)' }}>followers</p>
                    </div>
                    <div style={{ background: 'var(--color-background-secondary)', padding: '8px', cursor: 'pointer' }} onClick={() => setActiveTab(tabs.findIndex(o => o.label === 'following')) }>
                        <p style={{ fontSize:'18px', fontWeight:500, color:'#0F6E56' }}>31</p>
                        <p style={{ fontSize:'10px', color:'var(--color-text-tertiary)' }}>following</p>
                    </div>
                </div>
            </div>

            <div style={{ display: 'flex', gap:0, border:'0.5px solid var(--color-border-tertiary)', background:'var(--color-background-primary)' }}>
                {tabs.map((tab, idx) => (
                    <div
                        key={idx} // or label
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

            <CurrentTabComponent />
        </main>
    );
}
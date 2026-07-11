'use client'; // not good practice, try to remove it for all pages !

import Link from "next/link";
import { useState } from "react";

export default function Groups() {
  const [showCreateGroup, setShowCreateGroup] = useState(false);

  return (
    <main className="main">
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
        <p style={{ fontSize:'13px', fontWeight:500, color:'var(--color-text-primary)' }}>browse all groups</p>
        <button className="btn btn-p" onClick={ () => setShowCreateGroup(!showCreateGroup) } style={{ display: 'flex', alignItems: 'center', gap:'4px' }}><i className="ti ti-plus" style={{ fontSize:'12px' }} aria-hidden="true"></i> create group</button>
      </div>

      { showCreateGroup ? <div id="create-group-panel">
        <div className="card">
          <p style={{ fontSize:'11px', fontWeight:500, color:'var(--color-text-primary)', marginBottom:'10px' }}>new group</p>
          <div className="form-row"><span className="form-label">title *</span><input className="inp" placeholder="group name" /></div>
          <div className="form-row"><span className="form-label">description *</span><textarea className="inp" style={{ height:'56px', resize: 'none' }} placeholder="what is this group about?"></textarea></div>
          <div className="form-row"><span className="form-label">invite members</span><input className="inp" placeholder="search by name..." /></div>
          <div style={{ display: 'flex', gap:'6px', marginTop:'4px' }}>
            <button className="btn btn-p">create →</button>
            <button className="btn btn-g" onClick={ () => setShowCreateGroup(!showCreateGroup) }>cancel</button>
          </div>
        </div>
      </div> : <></>}

      <input className="inp" placeholder="search groups..." style={{ fontSize:'12px' }} />

      <div className="card">
        <div style={{ display: 'flex', alignItems: 'center', gap:'8px', marginBottom:'8px' }}>
          <span style={{ width:'8px', height:'8px', background:'#D4537E', flexShrink:0 }}></span>
          <p style={{ fontSize:'13px', fontWeight:500, color:'var(--color-text-primary)' }}>go devs</p>
          <span className="tag tag-pink" style={{ marginLeft: 'auto' }}>member</span>
        </div>
        <p style={{ fontSize:'12px', color:'var(--color-text-secondary)', marginBottom:'8px' }}>A space for Go enthusiasts — sharing projects, tips and co-building the social network backend.</p>
        <div style={{ fontSize:'11px', color:'var(--color-text-tertiary)', marginBottom:'10px' }}>34 members · created by Amir K.</div>
        <div className="event-card" style={{ marginBottom:'8px' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom:'4px' }}>
            <p style={{ fontSize:'12px', fontWeight:500, color:'var(--color-text-primary)' }}>Docker deep dive</p>
            <span className="tag tag-amber">Fri 18:00</span>
          </div>
          <p style={{ fontSize:'11px', color:'var(--color-text-secondary)', marginBottom:'6px' }}>Multi-stage builds, compose, and deployment pipelines.</p>
          <div style={{ display: 'flex', gap:'6px' }}><button className="btn btn-t" style={{ fontSize:'10px', display: 'flex', alignItems: 'center', gap:'3px' }}><i className="ti ti-check" style={{ fontSize:'11px' }} aria-hidden="true"></i> going (11)</button><button className="btn btn-g" style={{ fontSize:'10px' }}>not going (3)</button></div>
        </div>
        <div style={{ display: 'flex', gap:'6px' }}>
          <Link className="btn btn-p" href="/messages" style={{ fontSize:'11px', display: 'flex', alignItems: 'center', gap:'3px' }}><i className="ti ti-message" style={{ fontSize:'12px' }} aria-hidden="true"></i> group chat</Link>
          <button className="btn btn-g" style={{ fontSize:'11px', display: 'flex', alignItems: 'center', gap:'3px' }}><i className="ti ti-calendar-plus" style={{ fontSize:'12px' }} aria-hidden="true"></i> add event</button>
          <button className="btn btn-g" style={{ fontSize:'11px', display: 'flex', alignItems: 'center', gap:'3px' }}><i className="ti ti-user-plus" style={{ fontSize:'12px' }} aria-hidden="true"></i> invite</button>
        </div>
      </div>

      <div className="card">
        <div style={{ display: 'flex', alignItems: 'center', gap:'8px', marginBottom:'8px' }}>
          <span style={{ width:'8px', height:'8px', background:'#534AB7', flexShrink:0 }}></span>
          <p style={{ fontSize:'13px', fontWeight:500, color:'var(--color-text-primary)' }}>design systems</p>
          <span className="tag tag-gray" style={{ marginLeft: 'auto' }}>not a member</span>
        </div>
        <p style={{ fontSize:'12px', color:'var(--color-text-secondary)', marginBottom:'8px' }}>Discussing component libraries, tokens, and consistent UI patterns across projects.</p>
        <div style={{ fontSize:'11px', color:'var(--color-text-tertiary)', marginBottom:'10px' }}>21 members · created by Selin R.</div>
        <button className="btn btn-p" style={{ fontSize:'11px' }}>request to join →</button>
      </div>
    </main>
  );
}
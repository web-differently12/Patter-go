import React, { useState } from 'react';
import { OverviewPage } from './OverviewPage';
import { VoiceFailuresPage } from './VoiceFailuresPage';
import { CampaignsPage } from './CampaignsPage';
import { MeetingsPage } from './MeetingsPage';
import { BillingPage } from './BillingPage';

export type AnalyticsTab = 'overview' | 'failures' | 'campaigns' | 'meetings' | 'billing';

export const EchoFlowAnalytics: React.FC = () => {
  const [activeTab, setActiveTab] = useState<AnalyticsTab>('overview');

  return (
    <div className="min-h-screen bg-[#08090a] text-slate-100 font-sans">
      {/* Top Analytics Navigation Bar */}
      <div className="bg-slate-950/90 border-b border-slate-800 px-6 py-3 flex items-center justify-between sticky top-0 z-50 backdrop-blur-md">
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 rounded-lg bg-gradient-to-tr from-cyan-500 to-indigo-500 flex items-center justify-center font-bold text-slate-950 text-sm">
            EF
          </div>
          <span className="font-extrabold tracking-tight text-white text-base">EchoFlow Analytics</span>
        </div>

        <nav className="flex items-center gap-1 bg-slate-900 border border-slate-800 p-1 rounded-xl text-xs">
          <button
            onClick={() => setActiveTab('overview')}
            className={`px-3 py-1.5 rounded-lg transition-colors font-medium ${activeTab === 'overview' ? 'bg-cyan-500 text-slate-950 font-bold' : 'text-slate-400 hover:text-white'}`}
          >
            Vue d'Ensemble
          </button>
          <button
            onClick={() => setActiveTab('failures')}
            className={`px-3 py-1.5 rounded-lg transition-colors font-medium ${activeTab === 'failures' ? 'bg-cyan-500 text-slate-950 font-bold' : 'text-slate-400 hover:text-white'}`}
          >
            Diagnostic Échecs
          </button>
          <button
            onClick={() => setActiveTab('campaigns')}
            className={`px-3 py-1.5 rounded-lg transition-colors font-medium ${activeTab === 'campaigns' ? 'bg-cyan-500 text-slate-950 font-bold' : 'text-slate-400 hover:text-white'}`}
          >
            Campagnes Sortantes
          </button>
          <button
            onClick={() => setActiveTab('meetings')}
            className={`px-3 py-1.5 rounded-lg transition-colors font-medium ${activeTab === 'meetings' ? 'bg-cyan-500 text-slate-950 font-bold' : 'text-slate-400 hover:text-white'}`}
          >
            Réunions Visio
          </button>
          <button
            onClick={() => setActiveTab('billing')}
            className={`px-3 py-1.5 rounded-lg transition-colors font-medium ${activeTab === 'billing' ? 'bg-cyan-500 text-slate-950 font-bold' : 'text-slate-400 hover:text-white'}`}
          >
            Facturation & Wallet
          </button>
        </nav>
      </div>

      {/* Tab Content Rendering */}
      <main>
        {activeTab === 'overview' && <OverviewPage />}
        {activeTab === 'failures' && <VoiceFailuresPage />}
        {activeTab === 'campaigns' && <CampaignsPage />}
        {activeTab === 'meetings' && <MeetingsPage />}
        {activeTab === 'billing' && <BillingPage />}
      </main>
    </div>
  );
};

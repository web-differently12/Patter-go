import React, { useState } from 'react';
import type { ProviderMeta, ConnectedAccount, CustomAppFormValues, UserRole } from './types';
import { ProviderConnectModal } from './ProviderConnectModal';
import { ConnectedAccountBadge } from './ConnectedAccountBadge';

export const SUPPORTED_PROVIDERS: ProviderMeta[] = [
  {
    id: 'facebook',
    name: 'Facebook & Instagram',
    category: 'social_messaging',
    icon: 'FB',
    description: 'Messagerie directe Messenger, Instagram Direct & Meta Lead Ads.',
    developerConsoleUrl: 'https://developers.facebook.com/apps',
    supportsCustomApp: true,
  },
  {
    id: 'tiktok',
    name: 'TikTok Ads & Business',
    category: 'social_messaging',
    icon: 'TK',
    description: 'TikTok Messaging Direct & Lead Generation Campaigns.',
    developerConsoleUrl: 'https://developers.tiktok.com',
    supportsCustomApp: true,
  },
  {
    id: 'google',
    name: 'Google Workspace & Calendar',
    category: 'calendar_workspace',
    icon: 'GO',
    description: 'Google Calendar V2 sync, Gmail & Google Meet Bots.',
    developerConsoleUrl: 'https://console.cloud.google.com/apis/credentials',
    supportsCustomApp: true,
  },
  {
    id: 'slack',
    name: 'Slack Conversations',
    category: 'collaboration',
    icon: 'SL',
    description: 'Notifications temps reel et bots de reunion dans Slack.',
    developerConsoleUrl: 'https://api.slack.com/apps',
    supportsCustomApp: true,
  },
  {
    id: 'linkedin',
    name: 'LinkedIn Sales & Messaging',
    category: 'social_messaging',
    icon: 'LK',
    description: 'InMail automation & Sales Navigator Lead Capture.',
    developerConsoleUrl: 'https://www.linkedin.com/developers/apps',
    supportsCustomApp: true,
  },
  {
    id: 'hubspot',
    name: 'HubSpot CRM',
    category: 'ads_crm',
    icon: 'HS',
    description: 'Synchronisation automatique des contacts, opportunités et BANT scores.',
    developerConsoleUrl: 'https://developers.hubspot.com/docs/api/overview',
    supportsCustomApp: true,
  },
];

export const IntegrationsPage: React.FC = () => {
  const [currentRole, setRole] = useState<UserRole>('tenant_enterprise');
  const [selectedProvider, setSelectedProvider] = useState<ProviderMeta | null>(null);

  const [accounts, setAccounts] = useState<ConnectedAccount[]>([
    {
      id: 'acc_fb_01',
      tenantId: 'tenant_alpha_01',
      agencyId: 'agency_digital_scale',
      agencyName: 'Agence Digital Scale',
      provider: 'facebook',
      accountName: 'KallFlow Officiel (FB Page)',
      email: 'page@kallflow.ai',
      isCustomApp: false, // KallFlow Managed
      status: 'active',
      lastRefreshedAt: new Date().toISOString(),
    },
    {
      id: 'acc_tiktok_01',
      tenantId: 'tenant_beta_02',
      agencyId: 'agency_digital_scale',
      agencyName: 'Agence Digital Scale',
      provider: 'tiktok',
      accountName: 'KallFlow Enterprise TikTok App',
      email: 'tiktok-dev@kallflow.ai',
      isCustomApp: true, // Custom App BYO-App
      status: 'active',
      lastRefreshedAt: new Date().toISOString(),
    },
    {
      id: 'acc_google_01',
      tenantId: 'tenant_gamma_03',
      agencyId: 'agency_growth_lab',
      agencyName: 'Agence Growth Lab',
      provider: 'google',
      accountName: 'Google Calendar (Growth Lab)',
      email: 'calendar@growthlab.com',
      isCustomApp: true,
      status: 'active',
      lastRefreshedAt: new Date().toISOString(),
    },
  ]);

  // Filter visible accounts based on 3-tier hierarchy
  const visibleAccounts = accounts.filter((a) => {
    if (currentRole === 'super_admin') return true; // Super Admin sees all
    if (currentRole === 'agency') return a.agencyId === 'agency_digital_scale'; // Agency sees its sub-tenants
    return a.tenantId === 'tenant_alpha_01'; // Tenant Enterprise sees only its own
  });

  const handleConnectManaged = async () => {
    if (!selectedProvider) return;
    const newAccount: ConnectedAccount = {
      id: `acc_${selectedProvider.id}_${Date.now()}`,
      tenantId: 'tenant_alpha_01',
      agencyId: 'agency_digital_scale',
      agencyName: 'Agence Digital Scale',
      provider: selectedProvider.id,
      accountName: `${selectedProvider.name} (Managed)`,
      email: `user@kallflow.ai`,
      isCustomApp: false,
      status: 'active',
      lastRefreshedAt: new Date().toISOString(),
    };
    setAccounts((prev) => [...prev.filter((a) => a.provider !== selectedProvider.id), newAccount]);
  };

  const handleConnectCustomApp = async (values: CustomAppFormValues) => {
    if (!selectedProvider) return;
    const newAccount: ConnectedAccount = {
      id: `acc_${selectedProvider.id}_custom_${Date.now()}`,
      tenantId: 'tenant_alpha_01',
      agencyId: 'agency_digital_scale',
      agencyName: 'Agence Digital Scale',
      provider: selectedProvider.id,
      accountName: `${selectedProvider.name} (Custom App: ${values.clientId.slice(0, 8)}...)`,
      email: `dev-app@kallflow.ai`,
      isCustomApp: true,
      status: 'active',
      lastRefreshedAt: new Date().toISOString(),
    };
    setAccounts((prev) => [...prev.filter((a) => a.provider !== selectedProvider.id), newAccount]);
  };

  const handleTestConnection = async (id: string) => {
    await new Promise((res) => setTimeout(res, 800));
  };

  const handleRefreshToken = async (id: string) => {
    setAccounts((prev) =>
      prev.map((a) => (a.id === id ? { ...a, lastRefreshedAt: new Date().toISOString() } : a))
    );
  };

  const handleDisconnect = async (id: string) => {
    setAccounts((prev) => prev.filter((a) => a.id !== id));
  };

  return (
    <div className="space-y-6 text-slate-100 p-6 bg-[#08090a] min-h-screen font-sans">
      {/* Header & Role Switcher */}
      <div className="border-b border-zinc-800 pb-4 flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
        <div>
          <div className="flex items-center gap-2">
            <span className="text-xl font-black bg-clip-text text-transparent bg-gradient-to-r from-violet-400 to-cyan-400">KallFlow</span>
            <span className="text-xs bg-cyan-950 text-cyan-300 border border-cyan-800 px-2 py-0.5 rounded font-mono font-bold">Hybrid Integration Engine</span>
          </div>
          <h1 className="text-2xl font-bold tracking-tight text-white mt-1">Integrations Tierces & Meta / TikTok OAuth</h1>
          <p className="text-xs text-slate-400">Mode Managed 1-Click ou Application Developpeur Privee (BYO-App) avec gouvernance 3-tier.</p>
        </div>

        {/* 3-Tier Multi-Tenant Role Selector Simulation */}
        <div className="bg-zinc-950 border border-zinc-800 p-1.5 rounded-2xl flex items-center gap-1 text-xs">
          <span className="text-[10px] uppercase font-bold text-slate-500 px-2">Role:</span>
          <button
            onClick={() => setRole('tenant_enterprise')}
            className={`px-3 py-1.5 rounded-xl font-semibold transition-all ${
              currentRole === 'tenant_enterprise' ? 'bg-cyan-500 text-slate-950 font-bold' : 'text-slate-400 hover:text-white'
            }`}
          >
            Tenant Enterprise
          </button>
          <button
            onClick={() => setRole('agency')}
            className={`px-3 py-1.5 rounded-xl font-semibold transition-all ${
              currentRole === 'agency' ? 'bg-violet-600 text-white font-bold' : 'text-slate-400 hover:text-white'
            }`}
          >
            Agence
          </button>
          <button
            onClick={() => setRole('super_admin')}
            className={`px-3 py-1.5 rounded-xl font-semibold transition-all ${
              currentRole === 'super_admin' ? 'bg-rose-600 text-white font-bold' : 'text-slate-400 hover:text-white'
            }`}
          >
            Super Admin
          </button>
        </div>
      </div>

      {/* Hierarchy Info Badge */}
      <div className="p-3 bg-zinc-950 border border-zinc-800 rounded-xl text-xs flex justify-between items-center text-slate-400">
        <div>
          {currentRole === 'super_admin' && <span>[Super Admin Platform] Acces support global sur toutes les agences et sub-tenants ({visibleAccounts.length} comptes).</span>}
          {currentRole === 'agency' && <span>[Agence Digital Scale] Gestion des applications developpeur et sub-tenants rattaches ({visibleAccounts.length} comptes).</span>}
          {currentRole === 'tenant_enterprise' && <span>[Tenant Enterprise] Connexion directe Managed ou Application Privee anonyme Meta / TikTok.</span>}
        </div>
        <span className="text-[10px] text-cyan-400 font-mono">Nango Proxy Token Engine Active</span>
      </div>

      {/* Grid of Providers */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {SUPPORTED_PROVIDERS.map((provider) => {
          const connected = visibleAccounts.find((a) => a.provider === provider.id);

          return (
            <div key={provider.id} className="bg-zinc-950 border border-zinc-800 rounded-2xl p-5 shadow-xl flex flex-col justify-between space-y-4">
              <div>
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <span className="text-xs font-mono font-bold px-2 py-1 bg-zinc-800 text-zinc-300 rounded">[{provider.icon}]</span>
                    <div>
                      <h3 className="font-bold text-white text-base">{provider.name}</h3>
                      <span className="text-[10px] text-slate-500 uppercase tracking-wider font-semibold">{provider.category}</span>
                    </div>
                  </div>
                </div>
                <p className="text-xs text-slate-400 mt-3 leading-relaxed">{provider.description}</p>
              </div>

              {/* Connected Account Badge or Connect Trigger Button */}
              {connected ? (
                <div>
                  {currentRole !== 'tenant_enterprise' && connected.agencyName && (
                    <div className="text-[10px] text-violet-400 font-mono mb-1">Rattache a: {connected.agencyName} ({connected.tenantId})</div>
                  )}
                  <ConnectedAccountBadge
                    account={connected}
                    onTestConnection={handleTestConnection}
                    onRefreshToken={handleRefreshToken}
                    onDisconnect={handleDisconnect}
                  />
                </div>
              ) : (
                <button
                  onClick={() => setSelectedProvider(provider)}
                  className="w-full py-2.5 px-4 rounded-xl font-bold bg-gradient-to-r from-violet-600 to-cyan-500 hover:from-violet-500 hover:to-cyan-400 text-slate-950 shadow-lg transition-all text-xs flex items-center justify-center gap-2"
                >
                  <span>Connecter {provider.name}</span>
                  <span className="text-[10px] bg-slate-950/30 px-1.5 py-0.5 rounded text-white">Managed / BYO-App</span>
                </button>
              )}
            </div>
          );
        })}
      </div>

      {/* Provider Modal Trigger */}
      {selectedProvider && (
        <ProviderConnectModal
          provider={selectedProvider}
          isOpen={!!selectedProvider}
          onClose={() => setSelectedProvider(null)}
          onConnectManaged={handleConnectManaged}
          onConnectCustomApp={handleConnectCustomApp}
        />
      )}
    </div>
  );
};

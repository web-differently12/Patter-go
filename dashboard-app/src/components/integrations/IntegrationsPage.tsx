import React, { useState } from 'react';
import type { ProviderMeta, ConnectedAccount, CustomAppFormValues } from './types';
import { ProviderConnectModal } from './ProviderConnectModal';
import { ConnectedAccountBadge } from './ConnectedAccountBadge';

export const SUPPORTED_PROVIDERS: ProviderMeta[] = [
  {
    id: 'facebook',
    name: 'Facebook & Instagram',
    category: 'social_messaging',
    icon: '💬',
    description: 'Messagerie directe Messenger, Instagram Direct & Meta Lead Ads.',
    developerConsoleUrl: 'https://developers.facebook.com/apps',
    supportsCustomApp: true,
  },
  {
    id: 'tiktok',
    name: 'TikTok Ads & Business',
    category: 'social_messaging',
    icon: '🎵',
    description: 'TikTok Messaging Direct & Lead Generation Campaigns.',
    developerConsoleUrl: 'https://developers.tiktok.com',
    supportsCustomApp: true,
  },
  {
    id: 'google',
    name: 'Google Workspace & Calendar',
    category: 'calendar_workspace',
    icon: '📅',
    description: 'Google Calendar V2 sync, Gmail & Google Meet Bots.',
    developerConsoleUrl: 'https://console.cloud.google.com/apis/credentials',
    supportsCustomApp: true,
  },
  {
    id: 'slack',
    name: 'Slack Conversations',
    category: 'collaboration',
    icon: '🪟',
    description: 'Notifications temps réel et bots de réunion dans Slack.',
    developerConsoleUrl: 'https://api.slack.com/apps',
    supportsCustomApp: true,
  },
  {
    id: 'linkedin',
    name: 'LinkedIn Sales & Messaging',
    category: 'social_messaging',
    icon: '👔',
    description: 'InMail automation & Sales Navigator Lead Capture.',
    developerConsoleUrl: 'https://www.linkedin.com/developers/apps',
    supportsCustomApp: true,
  },
  {
    id: 'hubspot',
    name: 'HubSpot CRM',
    category: 'ads_crm',
    icon: '🟧',
    description: 'Synchronisation automatique des contacts, opportunités et BANT scores.',
    developerConsoleUrl: 'https://developers.hubspot.com/docs/api/overview',
    supportsCustomApp: true,
  },
];

export const IntegrationsPage: React.FC = () => {
  const [selectedProvider, setSelectedProvider] = useState<ProviderMeta | null>(null);
  const [accounts, setAccounts] = useState<ConnectedAccount[]>([
    {
      id: 'acc_fb_01',
      tenantId: 'tenant_kallflow_123',
      provider: 'facebook',
      accountName: 'KallFlow Offiiciel (FB Page)',
      email: 'page@kallflow.ai',
      isCustomApp: false, // Managed mode
      status: 'active',
      lastRefreshedAt: new Date().toISOString(),
    },
    {
      id: 'acc_tiktok_01',
      tenantId: 'tenant_kallflow_123',
      provider: 'tiktok',
      accountName: 'KallFlow Enterprise TikTok App',
      email: 'tiktok-dev@kallflow.ai',
      isCustomApp: true, // Custom App BYO-App
      status: 'active',
      lastRefreshedAt: new Date().toISOString(),
    },
  ]);

  const handleConnectManaged = async () => {
    if (!selectedProvider) return;
    // Simulate OAuth Popup & success callback
    const newAccount: ConnectedAccount = {
      id: `acc_${selectedProvider.id}_${Date.now()}`,
      tenantId: 'tenant_kallflow_123',
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
    // Simulate Custom App OAuth authorization flow
    const newAccount: ConnectedAccount = {
      id: `acc_${selectedProvider.id}_custom_${Date.now()}`,
      tenantId: 'tenant_kallflow_123',
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
      {/* Header */}
      <div className="border-b border-zinc-800 pb-4 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <div className="flex items-center gap-2">
            <span className="text-xl font-black bg-clip-text text-transparent bg-gradient-to-r from-violet-400 to-cyan-400">KallFlow</span>
            <span className="text-xs bg-cyan-950 text-cyan-300 border border-cyan-800 px-2 py-0.5 rounded font-mono font-bold">Hybrid Integration Selector</span>
          </div>
          <h1 className="text-2xl font-bold tracking-tight text-white mt-1">Intégrations Tierces & OAuth</h1>
          <p className="text-xs text-slate-400">Choisissez entre la connexion 1-Click KallFlow Managed ou votre propre Application Développeur (Marque Blanche BYO-App).</p>
        </div>
      </div>

      {/* Grid of Providers */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {SUPPORTED_PROVIDERS.map((provider) => {
          const connected = accounts.find((a) => a.provider === provider.id);

          return (
            <div key={provider.id} className="bg-zinc-950 border border-zinc-800 rounded-2xl p-5 shadow-xl flex flex-col justify-between space-y-4">
              <div>
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <span className="text-3xl">{provider.icon}</span>
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
                <ConnectedAccountBadge
                  account={connected}
                  onTestConnection={handleTestConnection}
                  onRefreshToken={handleRefreshToken}
                  onDisconnect={handleDisconnect}
                />
              ) : (
                <button
                  onClick={() => setSelectedProvider(provider)}
                  className="w-full py-2.5 px-4 rounded-xl font-bold bg-gradient-to-r from-violet-600 to-cyan-500 hover:from-violet-500 hover:to-cyan-400 text-slate-950 shadow-lg transition-all text-xs flex items-center justify-center gap-2"
                >
                  <span>Connecter {provider.name}</span>
                  <span className="text-[10px] bg-slate-950/30 px-1.5 py-0.5 rounded text-white">Managed / BYO-App →</span>
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

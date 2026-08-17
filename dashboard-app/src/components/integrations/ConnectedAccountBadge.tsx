import React, { useState } from 'react';
import type { ConnectedAccount } from './types';

export interface ConnectedAccountBadgeProps {
  account: ConnectedAccount;
  onTestConnection: (id: string) => Promise<void>;
  onRefreshToken: (id: string) => Promise<void>;
  onDisconnect: (id: string) => Promise<void>;
}

export const ConnectedAccountBadge: React.FC<ConnectedAccountBadgeProps> = ({
  account,
  onTestConnection,
  onRefreshToken,
  onDisconnect,
}) => {
  const [testing, setTesting] = useState(false);
  const [refreshing, setRefreshing] = useState(false);
  const [disconnecting, setDisconnecting] = useState(false);
  const [testSuccess, setTestSuccess] = useState<boolean | null>(null);

  const handleTest = async () => {
    setTesting(true);
    setTestSuccess(null);
    try {
      await onTestConnection(account.id);
      setTestSuccess(true);
      setTimeout(() => setTestSuccess(null), 3000);
    } catch {
      setTestSuccess(false);
    } finally {
      setTesting(false);
    }
  };

  const handleRefresh = async () => {
    setRefreshing(true);
    try {
      await onRefreshToken(account.id);
    } finally {
      setRefreshing(false);
    }
  };

  const handleDisconnect = async () => {
    if (!window.confirm(`Voulez-vous vraiment déconnecter le compte ${account.accountName} ?`)) return;
    setDisconnecting(true);
    try {
      await onDisconnect(account.id);
    } finally {
      setDisconnecting(false);
    }
  };

  return (
    <div className="p-4 bg-zinc-950 border border-zinc-800 rounded-2xl shadow-lg space-y-3 font-sans text-xs">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-full bg-zinc-800 flex items-center justify-center font-bold text-white overflow-hidden text-sm border border-zinc-700 font-mono">
            {account.accountAvatar ? (
              <img src={account.accountAvatar} alt={account.accountName} className="w-full h-full object-cover" />
            ) : (
              account.accountName.slice(0, 2).toUpperCase()
            )}
          </div>
          <div>
            <div className="flex items-center gap-2">
              <span className="font-bold text-white text-sm">{account.accountName}</span>
              <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" title="Actif"></span>
            </div>
            {account.email && <div className="text-slate-400 text-[11px]">{account.email}</div>}
          </div>
        </div>

        {/* Mode Badge */}
        <span
          className={`px-2.5 py-1 rounded-lg text-[10px] font-extrabold border ${
            account.isCustomApp
              ? 'bg-cyan-950/80 text-cyan-300 border-cyan-800/80'
              : 'bg-violet-950/80 text-violet-300 border-violet-800/80'
          }`}
        >
          {account.isCustomApp ? 'Application Privee (BYO-App)' : 'KallFlow Managed'}
        </span>
      </div>

      {testSuccess !== null && (
        <div className={`p-2 rounded-lg text-[11px] font-semibold border ${
          testSuccess ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/30' : 'bg-rose-500/10 text-rose-400 border-rose-500/30'
        }`}>
          {testSuccess ? 'Test de connexion reussi ! Token valide.' : 'Erreur lors du test de connexion.'}
        </div>
      )}

      {/* Action Buttons */}
      <div className="flex items-center justify-between pt-2 border-t border-zinc-800/80">
        <span className="text-[11px] text-slate-500">
          Dernier refresh: {new Date(account.lastRefreshedAt).toLocaleTimeString()}
        </span>

        <div className="flex gap-2">
          <button
            onClick={handleTest}
            disabled={testing}
            className="px-2.5 py-1 rounded-lg bg-zinc-900 hover:bg-zinc-800 text-slate-300 hover:text-white border border-zinc-800 transition-colors disabled:opacity-50"
          >
            {testing ? 'Test...' : 'Tester'}
          </button>
          <button
            onClick={handleRefresh}
            disabled={refreshing}
            className="px-2.5 py-1 rounded-lg bg-zinc-900 hover:bg-zinc-800 text-slate-300 hover:text-white border border-zinc-800 transition-colors disabled:opacity-50"
          >
            {refreshing ? 'Refresh...' : 'Actualiser Token'}
          </button>
          <button
            onClick={handleDisconnect}
            disabled={disconnecting}
            className="px-2.5 py-1 rounded-lg bg-rose-500/10 hover:bg-rose-500/20 text-rose-400 border border-rose-500/30 transition-colors disabled:opacity-50"
          >
            {disconnecting ? 'Suppression...' : 'Deconnecter'}
          </button>
        </div>
      </div>
    </div>
  );
};

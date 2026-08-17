import React, { useState } from 'react';
import type { IntegrationMode, ProviderMeta, CustomAppFormValues } from './types';

export interface ProviderConnectModalProps {
  provider: ProviderMeta;
  isOpen: boolean;
  onClose: () => void;
  onConnectManaged: () => Promise<void>;
  onConnectCustomApp: (values: CustomAppFormValues) => Promise<void>;
}

export const CALLBACK_URL = 'https://api.kallflow.ai/v1/integrations/oauth/callback';

export const ProviderConnectModal: React.FC<ProviderConnectModalProps> = ({
  provider,
  isOpen,
  onClose,
  onConnectManaged,
  onConnectCustomApp,
}) => {
  const [mode, setMode] = useState<IntegrationMode>('managed');
  const [clientId, setClientId] = useState('');
  const [clientSecret, setClientSecret] = useState('');
  const [scopes, setScopes] = useState('');
  const [showSecret, setShowSecret] = useState(false);
  const [showGuide, setShowGuide] = useState(false);
  const [copied, setCopied] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  if (!isOpen) return null;

  const handleCopyCallback = () => {
    navigator.clipboard.writeText(CALLBACK_URL);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const handleManagedSubmit = async () => {
    setLoading(true);
    setError(null);
    try {
      await onConnectManaged();
      onClose();
    } catch (e: any) {
      setError(e?.message || 'Echec de la connexion OAuth KallFlow Managed.');
    } finally {
      setLoading(false);
    }
  };

  const handleCustomSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!clientId.trim() || !clientSecret.trim()) {
      setError('Veuillez renseigner le Client ID et le Client Secret.');
      return;
    }
    setLoading(true);
    setError(null);
    try {
      await onConnectCustomApp({
        clientId: clientId.trim(),
        clientSecret: clientSecret.trim(),
        scopes: scopes.trim() || undefined,
      });
      onClose();
    } catch (e: any) {
      setError(e?.message || 'Erreur lors de la configuration de votre application privee.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm p-4 font-sans text-slate-100">
      <div className="w-full max-w-xl bg-[#08090a] border border-zinc-800 rounded-2xl shadow-2xl overflow-hidden flex flex-col max-h-[90vh]">
        {/* Header */}
        <div className="p-5 border-b border-zinc-800 flex items-center justify-between bg-zinc-950/60">
          <div className="flex items-center gap-3">
            <span className="text-xs font-mono px-2 py-1 bg-zinc-800 rounded text-zinc-300">[{provider.id.toUpperCase()}]</span>
            <div>
              <h2 className="text-lg font-bold text-white">Connecter {provider.name}</h2>
              <p className="text-xs text-slate-400">{provider.description}</p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="text-slate-400 hover:text-white text-lg px-2 py-1 rounded-lg hover:bg-zinc-800/60 transition-colors"
          >
            ✕
          </button>
        </div>

        {/* Mode Selector Tabs */}
        <div className="p-2 bg-zinc-950/80 border-b border-zinc-800 flex gap-2 text-xs">
          <button
            onClick={() => { setMode('managed'); setError(null); }}
            className={`flex-1 py-2.5 px-3 rounded-xl font-semibold transition-all flex items-center justify-center gap-2 ${
              mode === 'managed'
                ? 'bg-gradient-to-r from-violet-600 to-indigo-600 text-white shadow-md'
                : 'text-slate-400 hover:text-white hover:bg-zinc-900'
            }`}
          >
            <span>KallFlow Managed</span>
            <span className="text-[10px] bg-white/20 px-1.5 py-0.5 rounded font-bold">1-Click</span>
          </button>

          <button
            onClick={() => { setMode('custom_app'); setError(null); }}
            className={`flex-1 py-2.5 px-3 rounded-xl font-semibold transition-all flex items-center justify-center gap-2 ${
              mode === 'custom_app'
                ? 'bg-gradient-to-r from-cyan-500 to-blue-600 text-slate-950 font-extrabold shadow-md'
                : 'text-slate-400 hover:text-white hover:bg-zinc-900'
            }`}
          >
            <span>Application Privée (BYO-App)</span>
            <span className="text-[10px] bg-cyan-950 text-cyan-300 border border-cyan-800 px-1.5 py-0.5 rounded font-bold">Marque Blanche</span>
          </button>
        </div>

        {/* Content Body */}
        <div className="p-6 overflow-y-auto space-y-5 text-xs">
          {error && (
            <div className="p-3 bg-rose-500/10 border border-rose-500/30 rounded-xl text-rose-300 text-xs">
              [Erreur] {error}
            </div>
          )}

          {mode === 'managed' ? (
            <div className="space-y-4">
              <div className="p-4 bg-zinc-950 border border-zinc-800 rounded-xl space-y-2">
                <div className="font-bold text-slate-200 text-sm flex items-center gap-2">
                  <span className="text-emerald-400">[Valide]</span> Connexion Simplifiee Securisee
                </div>
                <p className="text-slate-400 leading-relaxed">
                  Connectez directement votre compte {provider.name} via l'application globale verifiee de KallFlow.
                  Aucune creation de compte developpeur requise. Les jetons sont chiffres et isoles sous votre <code className="text-cyan-400">tenant_id</code>.
                </p>
              </div>

              <div className="flex justify-end gap-3 pt-2">
                <button
                  onClick={onClose}
                  className="px-4 py-2 rounded-xl text-slate-400 hover:text-white bg-zinc-900 hover:bg-zinc-800 transition-colors"
                >
                  Annuler
                </button>
                <button
                  onClick={handleManagedSubmit}
                  disabled={loading}
                  className="px-5 py-2.5 rounded-xl font-bold bg-violet-600 hover:bg-violet-500 text-white shadow-lg transition-colors flex items-center gap-2 disabled:opacity-50"
                >
                  {loading ? 'Connexion en cours...' : `Connecter mon compte ${provider.name}`}
                </button>
              </div>
            </div>
          ) : (
            <form onSubmit={handleCustomSubmit} className="space-y-4">
              {/* Callback URL Copy Card */}
              <div className="p-3 bg-cyan-950/30 border border-cyan-800/50 rounded-xl space-y-2">
                <div className="flex justify-between items-center text-cyan-300 font-semibold">
                  <span>URL de Redirection Callback OAuth</span>
                  <button
                    type="button"
                    onClick={handleCopyCallback}
                    className="text-[11px] bg-cyan-500/20 hover:bg-cyan-500/30 text-cyan-300 px-2.5 py-1 rounded-lg border border-cyan-500/40 transition-colors"
                  >
                    {copied ? 'Copie !' : 'Copier URL'}
                  </button>
                </div>
                <code className="block bg-zinc-950 p-2 rounded text-[11px] text-slate-300 font-mono overflow-x-auto border border-zinc-800">
                  {CALLBACK_URL}
                </code>
              </div>

              {/* Developer Console Retractable Guide */}
              <div className="border border-zinc-800 rounded-xl overflow-hidden bg-zinc-950">
                <button
                  type="button"
                  onClick={() => setShowGuide(!showGuide)}
                  className="w-full p-3 text-left font-semibold text-slate-300 flex justify-between items-center hover:bg-zinc-900/60 transition-colors"
                >
                  <span>Guide pas-a-pas Console Developpeur {provider.name}</span>
                  <span>{showGuide ? '[-]' : '[+]'}</span>
                </button>
                {showGuide && (
                  <div className="p-4 border-t border-zinc-800 text-slate-400 space-y-2 bg-zinc-950/60 leading-relaxed">
                    <ol className="list-decimal list-inside space-y-1">
                      <li>Rendez-vous sur la <a href={provider.developerConsoleUrl} target="_blank" rel="noreferrer" className="text-cyan-400 hover:underline">Console Developpeur {provider.name}</a>.</li>
                      <li>Creez une nouvelle Application ou selectionnez votre App existante.</li>
                      <li>Dans les parametres OAuth, coller l'URL de redirection ci-dessus.</li>
                      <li>Copiez le <strong>Client ID / App ID</strong> et le <strong>Client Secret</strong> ci-dessous.</li>
                    </ol>
                  </div>
                )}
              </div>

              {/* Inputs */}
              <div className="space-y-3">
                <div>
                  <label className="block text-slate-300 font-semibold mb-1">Client ID / App ID</label>
                  <input
                    type="text"
                    value={clientId}
                    onChange={(e) => setClientId(e.target.value)}
                    placeholder={`ex: ${provider.id}_app_123456789`}
                    className="w-full bg-zinc-950 border border-zinc-800 rounded-xl px-3 py-2 text-slate-200 focus:outline-none focus:border-cyan-500 font-mono"
                  />
                </div>

                <div>
                  <label className="block text-slate-300 font-semibold mb-1">Client Secret / App Secret</label>
                  <div className="relative">
                    <input
                      type={showSecret ? 'text' : 'password'}
                      value={clientSecret}
                      onChange={(e) => setClientSecret(e.target.value)}
                      placeholder="••••••••••••••••••••••••••••"
                      className="w-full bg-zinc-950 border border-zinc-800 rounded-xl px-3 py-2 text-slate-200 focus:outline-none focus:border-cyan-500 font-mono pr-10"
                    />
                    <button
                      type="button"
                      onClick={() => setShowSecret(!showSecret)}
                      className="absolute right-3 top-2.5 text-slate-500 hover:text-slate-300 text-xs"
                    >
                      {showSecret ? 'Masquer' : 'Afficher'}
                    </button>
                  </div>
                </div>

                <div>
                  <label className="block text-slate-300 font-semibold mb-1">Scopes Personnalisés (Optionnel)</label>
                  <input
                    type="text"
                    value={scopes}
                    onChange={(e) => setScopes(e.target.value)}
                    placeholder="read,write,offline_access"
                    className="w-full bg-zinc-950 border border-zinc-800 rounded-xl px-3 py-2 text-slate-200 focus:outline-none focus:border-cyan-500 font-mono"
                  />
                </div>
              </div>

              {/* Actions */}
              <div className="flex justify-end gap-3 pt-3 border-t border-zinc-800">
                <button
                  type="button"
                  onClick={onClose}
                  className="px-4 py-2 rounded-xl text-slate-400 hover:text-white bg-zinc-900 hover:bg-zinc-800 transition-colors"
                >
                  Annuler
                </button>
                <button
                  type="submit"
                  disabled={loading}
                  className="px-5 py-2.5 rounded-xl font-bold bg-cyan-500 hover:bg-cyan-400 text-slate-950 shadow-lg transition-colors flex items-center gap-2 disabled:opacity-50"
                >
                  {loading ? 'Enregistrement...' : 'Enregistrer et Connecter mon App'}
                </button>
              </div>
            </form>
          )}
        </div>
      </div>
    </div>
  );
};

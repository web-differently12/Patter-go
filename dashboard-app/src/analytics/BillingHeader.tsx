import React from 'react';

export interface BillingHeaderProps {
  currency: 'EUR' | 'USD';
  setCurrency: (c: 'EUR' | 'USD') => void;
  theme: 'dark' | 'light';
  setTheme: (t: 'dark' | 'light') => void;
  onOpenTopUp: () => void;
}

export const BillingHeader: React.FC<BillingHeaderProps> = ({
  currency,
  setCurrency,
  theme,
  setTheme,
  onOpenTopUp,
}) => {
  return (
    <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 border-b border-zinc-800/80 pb-5 font-sans">
      <div>
        <div className="flex items-center gap-2.5">
          <h1 className="text-2xl font-semibold tracking-tight text-zinc-100">Facturation et Billing</h1>
          <span className="bg-emerald-950/40 text-emerald-300 border border-emerald-800/40 text-xs px-2.5 py-0.5 rounded-full flex items-center gap-1.5 font-medium">
            <span className="w-1.5 h-1.5 bg-emerald-500 rounded-full inline-block"></span>
            Active
          </span>
        </div>
        <p className="text-xs text-zinc-400 mt-1">Consommation en temps reel des agents KallFlow Engine, wallet et credits telecom.</p>
      </div>

      <div className="flex items-center gap-3">
        {/* Currency Switcher */}
        <div className="bg-zinc-900/80 border border-zinc-800/80 p-1 rounded-lg flex gap-1 text-xs">
          <button
            onClick={() => setCurrency('EUR')}
            className={`px-2.5 py-1 rounded font-medium transition-colors ${
              currency === 'EUR' ? 'bg-zinc-800 text-zinc-100' : 'text-zinc-400 hover:text-zinc-200'
            }`}
          >
            EUR
          </button>
          <button
            onClick={() => setCurrency('USD')}
            className={`px-2.5 py-1 rounded font-medium transition-colors ${
              currency === 'USD' ? 'bg-zinc-800 text-zinc-100' : 'text-zinc-400 hover:text-zinc-200'
            }`}
          >
            USD
          </button>
        </div>

        {/* Theme Mode Toggle (Dark / Light) */}
        <button
          onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
          className="p-2 bg-zinc-900/80 border border-zinc-800/80 hover:border-zinc-700/80 rounded-lg text-xs text-zinc-300 transition-colors font-mono"
        >
          {theme === 'dark' ? 'Mode Clair' : 'Mode Sombre'}
        </button>

        {/* Action TopUp Button - Linear/Vercel style */}
        <button
          onClick={onOpenTopUp}
          className="bg-white hover:bg-zinc-200 text-zinc-950 font-medium text-xs px-4 py-2 rounded-lg transition-colors shadow-sm"
        >
          + Recharger le Wallet
        </button>
      </div>
    </div>
  );
};

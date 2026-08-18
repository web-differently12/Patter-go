import React from 'react';

export interface WalletOverviewCardsProps {
  balanceEUR: number;
  currency: 'EUR' | 'USD';
  satoshis: number;
  estimatedDaysRemaining: number;
  agencyMarkupPercent: number;
}

export const WalletOverviewCards: React.FC<WalletOverviewCardsProps> = ({
  balanceEUR,
  currency,
  satoshis,
  estimatedDaysRemaining,
  agencyMarkupPercent,
}) => {
  const symbol = currency === 'EUR' ? 'EUR' : 'USD';
  const displayAmount = currency === 'EUR' ? balanceEUR : balanceEUR * 1.08;

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-5 font-sans">
      {/* Card 1: Wallet Balance & Satoshis */}
      <div className="bg-zinc-900/50 backdrop-blur-md border border-zinc-800/80 rounded-xl p-5 hover:border-zinc-700/80 transition-all duration-150 flex flex-col justify-between space-y-3">
        <div>
          <div className="text-xs font-medium text-zinc-400">Solde Wallet Unifie</div>
          <div className="text-2xl font-semibold tracking-tight text-zinc-100 font-mono mt-2">
            {displayAmount.toFixed(2)} {symbol}
          </div>
        </div>
        <div className="pt-2 border-t border-zinc-800/60 flex items-center justify-between text-xs">
          <span className="text-zinc-500">Equivalent Satoshis</span>
          <span className="text-amber-400/90 font-mono text-xs">{satoshis.toLocaleString()} sats</span>
        </div>
      </div>

      {/* Card 2: Burn Rate & Autonomy Progress */}
      <div className="bg-zinc-900/50 backdrop-blur-md border border-zinc-800/80 rounded-xl p-5 hover:border-zinc-700/80 transition-all duration-150 flex flex-col justify-between space-y-3">
        <div>
          <div className="flex justify-between items-center text-xs">
            <span className="font-medium text-zinc-400">Autonomie Estimee</span>
            <span className="text-zinc-300 font-mono text-xs">~{estimatedDaysRemaining} jours</span>
          </div>
          <div className="text-2xl font-semibold tracking-tight text-zinc-100 font-mono mt-2">
            18.40 {symbol} / jour
          </div>
        </div>
        {/* Fine 4px Burn Rate Progress Bar */}
        <div className="space-y-1.5 pt-2 border-t border-zinc-800/60">
          <div className="w-full bg-zinc-800 h-1 rounded-full overflow-hidden">
            <div className="bg-zinc-200 h-full rounded-full" style={{ width: `${Math.min(100, (estimatedDaysRemaining / 30) * 100)}%` }}></div>
          </div>
          <div className="text-[11px] text-zinc-500 flex justify-between">
            <span>Seuil alerte: 100.00 {symbol}</span>
            <span>Auto-topup: Actif</span>
          </div>
        </div>
      </div>

      {/* Card 3: Agency Markup Configuration */}
      <div className="bg-zinc-900/50 backdrop-blur-md border border-zinc-800/80 rounded-xl p-5 hover:border-zinc-700/80 transition-all duration-150 flex flex-col justify-between space-y-3">
        <div>
          <div className="flex justify-between items-center text-xs">
            <span className="font-medium text-zinc-400">Configuration Marge Agence</span>
            <span className="bg-zinc-800 text-zinc-300 border border-zinc-700/60 text-xs px-2 py-0.5 rounded-md font-mono">
              +{agencyMarkupPercent}%
            </span>
          </div>
          <div className="text-2xl font-semibold tracking-tight text-zinc-100 font-mono mt-2">
            Marque Blanche
          </div>
        </div>
        <p className="text-[11px] text-zinc-400 pt-2 border-t border-zinc-800/60 leading-relaxed">
          Marge appliquee automatiquement sur la consommation telecom et LLM de vos sub-tenants.
        </p>
      </div>
    </div>
  );
};

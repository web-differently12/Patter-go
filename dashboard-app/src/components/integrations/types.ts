export type UserRole = 'super_admin' | 'agency' | 'tenant_enterprise';

export interface AgencyPlan {
  id: string;
  name: string;
  maxSubTenants: number;
  whiteLabelEnabled: boolean;
  customAppAllowed: boolean;
  markupPercent: number;
}

export interface StripePlan {
  id: string;
  name: string;
  priceEUR: number;
  billingInterval: 'month' | 'year';
  stripePriceId: string;
  includedMinutes: number;
  includedTokens: number;
}

export interface TenantSubscription {
  tenantId: string;
  planId: string;
  planName: string;
  status: 'active' | 'past_due' | 'canceled';
  currentPeriodEnd: string;
  stripeCustomerId: string;
  stripeSubscriptionId: string;
  isAgencyPlan: boolean;
}

export interface TenantWallet {
  tenantId: string;
  creditsBalance: number; // 1 EUR = 1000 Credits
  autoReloadEnabled: boolean;
  autoReloadThreshold: number; // e.g. 10000
  autoReloadPackAmount: number; // e.g. 50000
}

export interface TenantQuotas {
  tenantId: string;
  stripeSubscriptionId: string;
  planTier: 'starter' | 'agency' | 'enterprise';
  whatsappMessagesUsedMonth: number;
  whatsappMessagesLimitMonth: number; // -1 for unlimited
  whatsappSessionsActive: number;
  whatsappSessionsLimit: number;
  mcpActiveClients: number;
  mcpClientsLimit: number; // -1 for unlimited
  billingCycleResetAt: string;
}

export interface CreditRate {
  service: string;
  unit: string;
  creditsPerUnit: number;
  note?: string;
}

export interface CreditPack {
  id: string;
  name: string;
  credits: number;
  priceEUR: number;
  stripePriceId: string;
}

export type IntegrationMode = 'managed' | 'custom_app';

export type ProviderId =
  | 'facebook'
  | 'tiktok'
  | 'google'
  | 'linkedin'
  | 'hubspot'
  | 'slack';

export interface ProviderMeta {
  id: ProviderId;
  name: string;
  category: 'social_messaging' | 'ads_crm' | 'calendar_workspace' | 'collaboration';
  icon: string;
  description: string;
  developerConsoleUrl: string;
  supportsCustomApp: boolean;
}

export interface TenantOAuthConfig {
  id: string;
  tenantId: string;
  provider: ProviderId;
  isCustomApp: boolean;
  customClientId?: string;
  customClientSecret?: string;
  customScopes?: string[];
  updatedAt: string;
}

export interface ConnectedAccount {
  id: string;
  tenantId: string;
  agencyId?: string;
  agencyName?: string;
  provider: ProviderId;
  accountName: string;
  accountAvatar?: string;
  email?: string;
  isCustomApp: boolean;
  status: 'active' | 'token_expired' | 'error';
  lastRefreshedAt: string;
}

export interface CustomAppFormValues {
  clientId: string;
  clientSecret: string;
  scopes?: string;
}

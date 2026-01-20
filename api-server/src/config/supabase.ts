import { createClient } from '@supabase/supabase-js';

const supabaseUrl = process.env.SUPABASE_URL!;
const supabaseServiceKey = process.env.SUPABASE_SERVICE_ROLE_KEY!;

export const supabaseAdmin = createClient(supabaseUrl, supabaseServiceKey, {
    auth: {
        autoRefreshToken: false,
        persistSession: false
    }
});

export const supabaseConfig = {
    url: supabaseUrl,
    anonKey: process.env.SUPABASE_ANON_KEY!,
    serviceRoleKey: supabaseServiceKey
};

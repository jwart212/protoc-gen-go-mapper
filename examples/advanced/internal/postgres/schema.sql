CREATE SCHEMA IF NOT EXISTS schm_pos;
SET search_path TO schm_pos;
-- Outbox table: stores domain events to be published to Kafka
CREATE TABLE IF NOT EXISTS schm_pos.outbox (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL,
  aggregate_id UUID,
  topic       TEXT NOT NULL,
  key         BYTEA,
  payload     BYTEA NOT NULL,
  headers     JSONB DEFAULT '{}'::jsonb,
  attempts    INT DEFAULT 0,
  status      ENUM('pending', 'publish','failed') NOT NULL DEFAULT 'pending', -- pending, sending, sent, failed, dead
  error_msg   TEXT NULL,
  retry_count BIGINT DEFAULT 0,
  retry_at TIMESTAMP NULL,
  created_at  TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at  TIMESTAMP WITH TIME ZONE DEFAULT now(),
  next_attempt_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_outbox_status_nextattempt ON outbox(status, next_attempt_at);



CREATE TABLE schm_pos.uoms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    symbol VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    deleted_by UUID
);

CREATE UNIQUE INDEX uq_uoms_tenant_code
ON schm_pos.uoms(tenant_id, code)
WHERE deleted_at IS NULL;

CREATE TABLE schm_pos.item_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    parent_id UUID,
    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    deleted_by UUID,
    CONSTRAINT fk_item_categories_parent
        FOREIGN KEY(parent_id)
        REFERENCES schm_pos.item_categories(id)
);

CREATE UNIQUE INDEX uq_item_categories_code
ON schm_pos.item_categories(tenant_id, code)
WHERE deleted_at IS NULL;

CREATE TYPE schm_pos.item_type AS ENUM (
    'RAW_MATERIAL',
    'SEMI_FINISHED',
    'FINISHED_GOOD',
    'PACKAGING',
    'SERVICE',
    'CONSUMABLE'
);

CREATE TABLE schm_pos.items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    category_id UUID NOT NULL,
    base_uom_id UUID NOT NULL,
    code VARCHAR(100) NOT NULL,
    sku VARCHAR(100),
    barcode VARCHAR(100),
    name VARCHAR(255) NOT NULL,
    item_type schm_pos.item_type NOT NULL,
    cost_price NUMERIC(18,4) DEFAULT 0,
    selling_price NUMERIC(18,4) DEFAULT 0,
    minimum_stock NUMERIC(18,4) DEFAULT 0,
    track_inventory BOOLEAN NOT NULL DEFAULT TRUE,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    deleted_by UUID,
    CONSTRAINT fk_items_category
        FOREIGN KEY(category_id)
        REFERENCES schm_pos.item_categories(id),

    CONSTRAINT fk_items_uom
        FOREIGN KEY(base_uom_id)
        REFERENCES schm_pos.uoms(id)
);

CREATE UNIQUE INDEX uq_items_code
ON schm_pos.items(tenant_id, code)
WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX uq_items_sku
ON schm_pos.items(tenant_id, sku)
WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX uq_items_barcode
ON schm_pos.items(tenant_id, barcode)
WHERE deleted_at IS NULL;


CREATE TABLE schm_pos.item_uoms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    item_id UUID NOT NULL,
    uom_id UUID NOT NULL,
    conversion_factor NUMERIC(18,6) NOT NULL,
    is_base BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_item_uom_item
        FOREIGN KEY(item_id)
        REFERENCES schm_pos.items(id),

    CONSTRAINT fk_item_uom_uom
        FOREIGN KEY(uom_id)
        REFERENCES schm_pos.uoms(id)
);


CREATE TABLE schm_pos.suppliers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    email VARCHAR(255),
    address TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX uq_suppliers_code
ON schm_pos.suppliers(tenant_id, code)
WHERE deleted_at IS NULL;


CREATE TABLE schm_pos.inventories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    item_id UUID NOT NULL,
    warehouse_id UUID NOT NULL,
    qty_on_hand NUMERIC(18,4) DEFAULT 0,
    qty_reserved NUMERIC(18,4) DEFAULT 0,
    qty_available NUMERIC(18,4) DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_inventory_item
        FOREIGN KEY(item_id)
        REFERENCES schm_pos.items(id)
);

CREATE UNIQUE INDEX uq_inventory_item_wh
ON schm_pos.inventories(
    tenant_id,
    warehouse_id,
    item_id
);

CREATE TYPE schm_pos.inventory_transaction_type AS ENUM (
    'PURCHASE',
    'SALE',
    'PRODUCTION_USAGE',
    'PRODUCTION_RESULT',
    'RETURN',
    'ADJUSTMENT',
    'WASTE',
    'TRANSFER'
);

CREATE TABLE schm_pos.inventory_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    item_id UUID NOT NULL,
    warehouse_id UUID NOT NULL,
    transaction_type schm_pos.inventory_transaction_type NOT NULL,
    reference_type VARCHAR(100),
    reference_id UUID,
    qty_in NUMERIC(18,4) DEFAULT 0,
    qty_out NUMERIC(18,4) DEFAULT 0,
    balance_after NUMERIC(18,4) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE schm_pos.boms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    product_id UUID NOT NULL,
    version VARCHAR(50),
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TYPE schm_pos.production_order_status AS ENUM (
    'DRAFT',
    'RELEASED',
    'IN_PROGRESS',
    'COMPLETED',
    'CANCELLED'
);


CREATE TABLE schm_pos.production_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    branch_id UUID NOT NULL,
    warehouse_id UUID NOT NULL,
    bom_id UUID NOT NULL,
    production_no VARCHAR(100) NOT NULL,
    production_date DATE NOT NULL,
    status schm_pos.production_order_status NOT NULL DEFAULT 'DRAFT',
    qty_target NUMERIC(18,4) NOT NULL,
    qty_actual NUMERIC(18,4),
    expected_yield NUMERIC(18,4),
    actual_yield NUMERIC(18,4),
    waste_qty NUMERIC(18,4),
    notes TEXT,
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_po_bom
        FOREIGN KEY(bom_id)
        REFERENCES schm_pos.boms(id),
    CONSTRAINT fk_po_branch
        FOREIGN KEY(branch_id)
        REFERENCES schm_pos.branches(id),
    CONSTRAINT fk_po_warehouse
        FOREIGN KEY(warehouse_id)
        REFERENCES schm_pos.warehouses(id)
);


CREATE TABLE schm_pos.production_order_materials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    production_order_id UUID NOT NULL,
    item_id UUID NOT NULL,
    uom_id UUID NOT NULL,
    batch_id UUID,
    qty_planned NUMERIC(18,4) NOT NULL,
    qty_actual NUMERIC(18,4),
    unit_cost NUMERIC(18,4),
    total_cost NUMERIC(18,4),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_pom_po
        FOREIGN KEY(production_order_id)
        REFERENCES schm_pos.production_orders(id),
    CONSTRAINT fk_pom_item
        FOREIGN KEY(item_id)
        REFERENCES schm_pos.items(id),
    CONSTRAINT fk_pom_uom
        FOREIGN KEY(uom_id)
        REFERENCES schm_pos.uoms(id),
    CONSTRAINT fk_pom_batch
        FOREIGN KEY(batch_id)
        REFERENCES schm_pos.inventory_batches(id)
);


CREATE TABLE schm_pos.production_order_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    production_order_id UUID NOT NULL,
    item_id UUID NOT NULL,
    uom_id UUID NOT NULL,
    batch_id UUID,
    qty NUMERIC(18,4) NOT NULL,
    unit_cost NUMERIC(18,4),
    total_cost NUMERIC(18,4),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_por_po
        FOREIGN KEY(production_order_id)
        REFERENCES schm_pos.production_orders(id),
    CONSTRAINT fk_por_item
        FOREIGN KEY(item_id)
        REFERENCES schm_pos.items(id),
    CONSTRAINT fk_por_uom
        FOREIGN KEY(uom_id)
        REFERENCES schm_pos.uoms(id)
);


CREATE TYPE schm_pos.purchase_status AS ENUM (
    'DRAFT',
    'CONFIRMED',
    'RECEIVED',
    'CANCELLED',
    'RETURE'
);


CREATE TABLE schm_pos.purchases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    supplier_id UUID NOT NULL,
    purchase_no VARCHAR(100) NOT NULL,
    purchase_date DATE NOT NULL,
    total_amount NUMERIC(18,2) DEFAULT 0,
    status schm_pos.purchase_status NOT NULL DEFAULT 'DRAFT',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TYPE schm_pos.sales_status AS ENUM (
    'DRAFT',
    'PAID',
    'PARTIAL',
    'VOID'
);


CREATE TABLE schm_pos.sales (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    customer_id UUID,
    invoice_no VARCHAR(100) NOT NULL,
    subtotal NUMERIC(18,2) DEFAULT 0,
    discount NUMERIC(18,2) DEFAULT 0,
    tax NUMERIC(18,2) DEFAULT 0,
    total NUMERIC(18,2) DEFAULT 0,
    status schm_pos.sales_status NOT NULL DEFAULT 'DRAFT',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);


CREATE TABLE schm_pos.branches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    email VARCHAR(255),
    address TEXT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    deleted_by UUID
);

CREATE UNIQUE INDEX uq_branches_tenant_code
ON schm_pos.branches(tenant_id, code)
WHERE deleted_at IS NULL;

CREATE TABLE schm_pos.warehouses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    branch_id UUID NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    warehouse_type schm_pos.warehouse_type NOT NULL DEFAULT 'MAIN',
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_warehouse_branch
        FOREIGN KEY(branch_id)
        REFERENCES schm_pos.branches(id)
);


CREATE TYPE schm_pos.warehouse_type AS ENUM (
    'MAIN',
    'KITCHEN',
    'OUTLET',
    'BAR',
    'PRODUCTION',
    'TRANSIT'
);

CREATE TABLE schm_pos.customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    code VARCHAR(50),
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    email VARCHAR(255),
    birth_date DATE,
    loyalty_points BIGINT DEFAULT 0,
    member_since DATE,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);



CREATE TABLE schm_pos.stock_transfers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    transfer_no VARCHAR(100) NOT NULL,
    source_warehouse_id UUID NOT NULL,
    destination_warehouse_id UUID NOT NULL,
    transfer_date TIMESTAMPTZ NOT NULL,
    status VARCHAR(50) NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE schm_pos.stock_transfer_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transfer_id UUID NOT NULL,
    item_id UUID NOT NULL,
    uom_id UUID NOT NULL,
    qty NUMERIC(18,4) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE schm_pos.stock_opnames (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    warehouse_id UUID NOT NULL,
    opname_no VARCHAR(100) NOT NULL,
    opname_date TIMESTAMPTZ NOT NULL,
    status VARCHAR(50) NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE schm_pos.stock_opname_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    opname_id UUID NOT NULL,
    item_id UUID NOT NULL,
    system_qty NUMERIC(18,4),
    actual_qty NUMERIC(18,4),
    variance_qty NUMERIC(18,4),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE schm_pos.inventory_batches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    warehouse_id UUID NOT NULL,
    item_id UUID NOT NULL,
    batch_no VARCHAR(100) NOT NULL,
    manufacture_date DATE,
    expired_date DATE,
    qty NUMERIC(18,4) NOT NULL,
    cost_per_unit NUMERIC(18,4),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE schm_pos.inventory_wastes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    warehouse_id UUID NOT NULL,
    waste_no VARCHAR(100) NOT NULL,
    waste_date TIMESTAMPTZ NOT NULL,
    reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE schm_pos.inventory_waste_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    waste_id UUID NOT NULL,
    item_id UUID NOT NULL,
    qty NUMERIC(18,4) NOT NULL,
    cost NUMERIC(18,4),
    notes TEXT
);

CREATE TABLE schm_pos.cash_shifts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    branch_id UUID NOT NULL,
    shift_no VARCHAR(100) NOT NULL,
    cashier_id UUID NOT NULL,
    opened_at TIMESTAMPTZ NOT NULL,
    closed_at TIMESTAMPTZ,
    opening_balance NUMERIC(18,2),
    closing_balance NUMERIC(18,2),
    expected_balance NUMERIC(18,2),
    variance NUMERIC(18,2),
    status VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE schm_pos.cash_movements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    cash_shift_id UUID NOT NULL,
    movement_type VARCHAR(50),
    amount NUMERIC(18,2),
    reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE schm_pos.kitchen_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    sale_id UUID NOT NULL,
    kitchen_no VARCHAR(100),
    status VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE schm_pos.kitchen_order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kitchen_order_id UUID NOT NULL,
    sale_item_id UUID NOT NULL,
    status VARCHAR(50),
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ
);

CREATE TABLE schm_pos.menu_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    code VARCHAR(50),
    name VARCHAR(255),
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE schm_pos.menus (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    category_id UUID NOT NULL,
    item_id UUID NOT NULL,
    code VARCHAR(50),
    name VARCHAR(255),
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE schm_pos.menu_prices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    menu_id UUID NOT NULL,
    branch_id UUID,
    price NUMERIC(18,2),
    effective_from TIMESTAMPTZ,
    effective_until TIMESTAMPTZ
);

ALTER TABLE schm_pos.inventories
ADD COLUMN branch_id UUID NOT NULL;


ALTER TABLE schm_pos.inventory_transactions
ADD COLUMN branch_id UUID NOT NULL;

ALTER TABLE schm_pos.inventory_transactions
ADD COLUMN batch_id UUID;


ALTER TABLE schm_pos.sales
ADD COLUMN branch_id UUID NOT NULL;

ALTER TABLE schm_pos.sales
ADD COLUMN cash_shift_id UUID;

ALTER TABLE schm_pos.sales
ADD COLUMN kitchen_order_id UUID;


ALTER TABLE schm_pos.purchases
ADD COLUMN branch_id UUID NOT NULL;

ALTER TABLE schm_pos.purchases
ADD COLUMN warehouse_id UUID NOT NULL;



ALTER TABLE schm_pos.production_orders
ADD COLUMN loss_qty NUMERIC(18,4);


ALTER TABLE schm_pos.boms
ADD COLUMN effective_from DATE;

ALTER TABLE schm_pos.boms
ADD COLUMN effective_until DATE;


CREATE INDEX idx_sales_tenant_branch
ON schm_pos.sales(tenant_id, branch_id);

CREATE INDEX idx_inventory_tenant_wh_item
ON schm_pos.inventories(
    tenant_id,
    warehouse_id,
    item_id
);

CREATE INDEX idx_batch_expired
ON schm_pos.inventory_batches(
    tenant_id,
    expired_date
);

CREATE INDEX idx_kitchen_status
ON schm_pos.kitchen_orders(
    tenant_id,
    status
);

CREATE INDEX idx_cash_shift_status
ON schm_pos.cash_shifts(
    tenant_id,
    status
);
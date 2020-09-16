/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */
package dm

import (
	"crypto/tls"
	"dmdb_exporter/dm/security"
	"dmdb_exporter/dm/util"
	"net"
	"strconv"
	"time"
)

const (
	Dm_build_1275 = 8192
	Dm_build_1276 = 2 * time.Second
)

type dm_build_1277 struct {
	dm_build_1278 *net.TCPConn
	dm_build_1279 *tls.Conn
	dm_build_1280 *util.Dm_build_912
	dm_build_1281 *DmConnection
	dm_build_1282 security.Cipher
	dm_build_1283 bool
	dm_build_1284 bool
	dm_build_1285 *security.DhKey
	dm_build_1286 string
	dm_build_1287 bool
}

func dm_build_1288(dm_build_1289 *DmConnection) (*dm_build_1277, error) {
	dm_build_1290, dm_build_1291 := dm_build_1293(dm_build_1289.dmConnector.host+":"+strconv.Itoa(dm_build_1289.dmConnector.port), time.Duration(dm_build_1289.dmConnector.socketTimeout)*time.Second)
	if dm_build_1291 != nil {
		return nil, dm_build_1291
	}

	dm_build_1292 := dm_build_1277{}
	dm_build_1292.dm_build_1278 = dm_build_1290
	dm_build_1292.dm_build_1280 = util.Dm_build_918(Dm_build_6, true)
	dm_build_1292.dm_build_1281 = dm_build_1289
	dm_build_1292.dm_build_1283 = false
	dm_build_1292.dm_build_1284 = false
	dm_build_1292.dm_build_1286 = ""
	dm_build_1292.dm_build_1287 = false
	dm_build_1289.Access = &dm_build_1292

	return &dm_build_1292, nil
}

func dm_build_1293(dm_build_1294 string, dm_build_1295 time.Duration) (*net.TCPConn, error) {
	dm_build_1296, dm_build_1297 := net.DialTimeout("tcp", dm_build_1294, dm_build_1295)
	if dm_build_1297 != nil {
		return nil, ECGO_COMMUNITION_ERROR.addDetail("\tdial address: " + dm_build_1294).throw()
	}

	if tcpConn, ok := dm_build_1296.(*net.TCPConn); ok {

		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(Dm_build_1276)
		tcpConn.SetNoDelay(true)

		tcpConn.SetReadBuffer(Dm_build_1275)
		tcpConn.SetWriteBuffer(Dm_build_1275)
		return tcpConn, nil
	}

	return nil, nil
}

func (dm_build_1299 *dm_build_1277) dm_build_1298(dm_build_1300 dm_build_123) bool {
	var dm_build_1301 = dm_build_1299.dm_build_1281.dmConnector.compress
	if dm_build_1300.dm_build_137() == Dm_build_33 || dm_build_1301 == Dm_build_81 {
		return false
	}

	if dm_build_1301 == Dm_build_79 {
		return true
	} else if dm_build_1301 == Dm_build_80 {
		return !dm_build_1299.dm_build_1281.Local && dm_build_1300.dm_build_135() > Dm_build_78
	}

	return false
}

func (dm_build_1303 *dm_build_1277) dm_build_1302(dm_build_1304 dm_build_123) bool {
	var dm_build_1305 = dm_build_1303.dm_build_1281.dmConnector.compress
	if dm_build_1304.dm_build_137() == Dm_build_33 || dm_build_1305 == Dm_build_81 {
		return false
	}

	if dm_build_1305 == Dm_build_79 {
		return true
	} else if dm_build_1305 == Dm_build_80 {
		return dm_build_1303.dm_build_1280.Dm_build_1210(Dm_build_41) == 1
	}

	return false
}

func (dm_build_1307 *dm_build_1277) dm_build_1306(dm_build_1308 dm_build_123) (err error) {
	defer func() {
		if p := recover(); p != nil {
			if _, ok := p.(string); ok {
				err = ECGO_COMMUNITION_ERROR.addDetail("\t" + p.(string)).throw()
			} else {
				panic(p)
			}
		}
	}()

	dm_build_1310 := dm_build_1308.dm_build_135()

	if dm_build_1310 > 0 {

		if dm_build_1307.dm_build_1298(dm_build_1308) {
			var retBytes, err = Compress(dm_build_1307.dm_build_1280, Dm_build_34, int(dm_build_1310), int(dm_build_1307.dm_build_1281.dmConnector.compressID))
			if err != nil {
				return err
			}

			dm_build_1307.dm_build_1280.Dm_build_971(Dm_build_34)

			dm_build_1307.dm_build_1280.Dm_build_1023(dm_build_1310)

			dm_build_1307.dm_build_1280.Dm_build_1044(retBytes)

			dm_build_1308.dm_build_136(int32(len(retBytes)) + ULINT_SIZE)

			dm_build_1307.dm_build_1280.Dm_build_1132(Dm_build_41, 1)
		}

		dm_build_1310 = dm_build_1308.dm_build_135()
		if dm_build_1307.dm_build_1284 {
			var retBytes = dm_build_1307.dm_build_1282.Encrypt(dm_build_1307.dm_build_1280.Dm_build_1237(Dm_build_34, make([]byte, dm_build_1310)), true)

			dm_build_1307.dm_build_1280.Dm_build_971(Dm_build_34)

			dm_build_1307.dm_build_1280.Dm_build_1044(retBytes)

			dm_build_1308.dm_build_136(int32(len(retBytes)))
		}
	}

	dm_build_1308.dm_build_132()
	if dm_build_1307.dm_build_1501(dm_build_1308) {
		if dm_build_1307.dm_build_1279 != nil {
			dm_build_1307.dm_build_1280.Dm_build_976(0)
			dm_build_1307.dm_build_1280.Dm_build_1001(dm_build_1307.dm_build_1279, false)
		}
	} else {
		dm_build_1307.dm_build_1280.Dm_build_976(0)
		dm_build_1307.dm_build_1280.Dm_build_1001(dm_build_1307.dm_build_1278, false)
	}
	return nil
}

func (dm_build_1312 *dm_build_1277) dm_build_1311(dm_build_1313 dm_build_123) (err error) {
	defer func() {
		if p := recover(); p != nil {
			if _, ok := p.(string); ok {
				err = ECGO_COMMUNITION_ERROR.addDetail("\t" + p.(string)).throw()
			} else {
				panic(p)
			}
		}
	}()

	dm_build_1315 := int32(0)
	if dm_build_1312.dm_build_1501(dm_build_1313) {
		if dm_build_1312.dm_build_1279 != nil {
			dm_build_1312.dm_build_1280.Dm_build_971(0)
			dm_build_1312.dm_build_1280.Dm_build_996(dm_build_1312.dm_build_1279, Dm_build_34)
			dm_build_1315 = dm_build_1313.dm_build_135()
			if dm_build_1315 > 0 {
				dm_build_1312.dm_build_1280.Dm_build_996(dm_build_1312.dm_build_1279, int(dm_build_1315))
			}
		}
	} else {

		dm_build_1312.dm_build_1280.Dm_build_971(0)
		dm_build_1312.dm_build_1280.Dm_build_996(dm_build_1312.dm_build_1278, Dm_build_34)
		dm_build_1315 = dm_build_1313.dm_build_135()

		if dm_build_1315 > 0 {
			dm_build_1312.dm_build_1280.Dm_build_996(dm_build_1312.dm_build_1278, int(dm_build_1315))
		}
	}

	dm_build_1313.dm_build_133()

	if dm_build_1315 <= 0 {
		return nil
	}

	if dm_build_1312.dm_build_1284 {
		ebytes := dm_build_1312.dm_build_1280.Dm_build_1237(Dm_build_34, make([]byte, dm_build_1315))
		bytes, err := dm_build_1312.dm_build_1282.Decrypt(ebytes, true)
		if err != nil {
			return err
		}
		dm_build_1312.dm_build_1280.Dm_build_971(Dm_build_34)
		dm_build_1312.dm_build_1280.Dm_build_1044(bytes)
		dm_build_1313.dm_build_136(int32(len(bytes)))
	}

	if dm_build_1312.dm_build_1302(dm_build_1313) {

		cbytes := dm_build_1312.dm_build_1280.Dm_build_1237(Dm_build_34+ULINT_SIZE, make([]byte, dm_build_1315-ULINT_SIZE))
		bytes, err := UnCompress(cbytes, int(dm_build_1312.dm_build_1281.dmConnector.compressID))
		if err != nil {
			return err
		}
		dm_build_1312.dm_build_1280.Dm_build_971(Dm_build_34)
		dm_build_1312.dm_build_1280.Dm_build_1044(bytes)
		dm_build_1313.dm_build_136(int32(len(bytes)))
	}
	return nil
}

func (dm_build_1317 *dm_build_1277) dm_build_1316(dm_build_1318 dm_build_123) (dm_build_1319 interface{}, dm_build_1320 error) {
	dm_build_1320 = dm_build_1318.dm_build_127(dm_build_1318)
	if dm_build_1320 != nil {
		return nil, dm_build_1320
	}

	dm_build_1320 = dm_build_1317.dm_build_1306(dm_build_1318)
	if dm_build_1320 != nil {
		return nil, dm_build_1320
	}

	dm_build_1320 = dm_build_1317.dm_build_1311(dm_build_1318)
	if dm_build_1320 != nil {
		return nil, dm_build_1320
	}

	return dm_build_1318.dm_build_131(dm_build_1318)
}

func (dm_build_1322 *dm_build_1277) dm_build_1321() (*dm_build_528, error) {

	Dm_build_1323 := dm_build_534(dm_build_1322)
	_, dm_build_1324 := dm_build_1322.dm_build_1316(Dm_build_1323)
	if dm_build_1324 != nil {
		return nil, dm_build_1324
	}

	return Dm_build_1323, nil
}

func (dm_build_1326 *dm_build_1277) dm_build_1325() error {

	dm_build_1327 := dm_build_405(dm_build_1326)
	_, dm_build_1328 := dm_build_1326.dm_build_1316(dm_build_1327)
	if dm_build_1328 != nil {
		return dm_build_1328
	}

	return nil
}

func (dm_build_1330 *dm_build_1277) dm_build_1329() error {

	var dm_build_1331 *dm_build_528
	var err error
	if dm_build_1331, err = dm_build_1330.dm_build_1321(); err != nil {
		return err
	}

	if dm_build_1330.dm_build_1281.sslEncrypt == 2 {
		if err = dm_build_1330.dm_build_1497(false); err != nil {
			return ECGO_INIT_SSL_FAILED.addDetail("\n" + err.Error()).throw()
		}
	} else if dm_build_1330.dm_build_1281.sslEncrypt == 1 {
		if err = dm_build_1330.dm_build_1497(true); err != nil {
			return ECGO_INIT_SSL_FAILED.addDetail("\n" + err.Error()).throw()
		}
	}

	if dm_build_1330.dm_build_1284 || dm_build_1330.dm_build_1283 {
		k, err := dm_build_1330.dm_build_1487()
		if err != nil {
			return err
		}
		sessionKey := security.ComputeSessionKey(k, dm_build_1331.Dm_build_532)
		encryptType := dm_build_1331.dm_build_530
		hashType := int(dm_build_1331.Dm_build_531)
		if encryptType == -1 {
			encryptType = security.DES_CFB
		}
		if hashType == -1 {
			hashType = security.MD5
		}
		err = dm_build_1330.dm_build_1490(encryptType, sessionKey, dm_build_1330.dm_build_1281.dmConnector.cipherPath, hashType)
		if err != nil {
			return err
		}
	}

	if err := dm_build_1330.dm_build_1325(); err != nil {
		return err
	}
	return nil
}

func (dm_build_1334 *dm_build_1277) Dm_build_1333(dm_build_1335 *DmStatement) error {
	dm_build_1336 := dm_build_557(dm_build_1334, dm_build_1335)
	_, dm_build_1337 := dm_build_1334.dm_build_1316(dm_build_1336)
	if dm_build_1337 != nil {
		return dm_build_1337
	}

	return nil
}

func (dm_build_1339 *dm_build_1277) Dm_build_1338(dm_build_1340 int32) error {
	dm_build_1341 := dm_build_567(dm_build_1339, dm_build_1340)
	_, dm_build_1342 := dm_build_1339.dm_build_1316(dm_build_1341)
	if dm_build_1342 != nil {
		return dm_build_1342
	}

	return nil
}

func (dm_build_1344 *dm_build_1277) Dm_build_1343(dm_build_1345 *DmStatement, dm_build_1346 bool, dm_build_1347 int16) (*execInfo, error) {
	dm_build_1348 := dm_build_441(dm_build_1344, dm_build_1345, dm_build_1346, dm_build_1347)
	dm_build_1349, dm_build_1350 := dm_build_1344.dm_build_1316(dm_build_1348)
	if dm_build_1350 != nil {
		return nil, dm_build_1350
	}
	return dm_build_1349.(*execInfo), nil
}

func (dm_build_1352 *dm_build_1277) Dm_build_1351(dm_build_1353 *DmStatement, dm_build_1354 int16) (*execInfo, error) {
	return dm_build_1352.Dm_build_1343(dm_build_1353, false, Dm_build_85)
}

func (dm_build_1356 *dm_build_1277) Dm_build_1355(dm_build_1357 *DmStatement, dm_build_1358 []OptParameter) (*execInfo, error) {
	dm_build_1359, dm_build_1360 := dm_build_1356.dm_build_1316(dm_build_216(dm_build_1356, dm_build_1357, dm_build_1358))
	if dm_build_1360 != nil {
		return nil, dm_build_1360
	}

	return dm_build_1359.(*execInfo), nil
}

func (dm_build_1362 *dm_build_1277) Dm_build_1361(dm_build_1363 *DmStatement, dm_build_1364 int16) (*execInfo, error) {
	return dm_build_1362.Dm_build_1343(dm_build_1363, true, dm_build_1364)
}

func (dm_build_1366 *dm_build_1277) Dm_build_1365(dm_build_1367 *DmStatement, dm_build_1368 [][]interface{}) (*execInfo, error) {
	var dm_build_1369 = false

	var dm_build_1370 = make([]interface{}, dm_build_1367.paramCount)
	for icol := 0; icol < int(dm_build_1367.paramCount); icol++ {
		if dm_build_1367.params[icol].ioType == IO_TYPE_OUT {
			continue
		}
		if dm_build_1366.dm_build_1470(dm_build_1370, dm_build_1368[0], icol) {

			if !dm_build_1369 {
				preExecute := dm_build_431(dm_build_1366, dm_build_1367, dm_build_1367.params)
				dm_build_1366.dm_build_1316(preExecute)
				dm_build_1369 = true
			}

			dm_build_1366.dm_build_1476(dm_build_1367, dm_build_1367.params[icol], icol, dm_build_1368[0][icol].(iOffRowBinder))
			dm_build_1370[icol] = ParamDataEnum_OFF_ROW
		}
	}

	var dm_build_1371 = make([][]interface{}, 1, 1)
	dm_build_1371[0] = dm_build_1370

	dm_build_1372 := dm_build_239(dm_build_1366, dm_build_1367, dm_build_1371)
	dm_build_1373, dm_build_1374 := dm_build_1366.dm_build_1316(dm_build_1372)
	if dm_build_1374 != nil {
		return nil, dm_build_1374
	}
	return dm_build_1373.(*execInfo), nil
}

func (dm_build_1376 *dm_build_1277) Dm_build_1375(dm_build_1377 *DmStatement, dm_build_1378 int16) (*execInfo, error) {
	dm_build_1379 := dm_build_419(dm_build_1376, dm_build_1377, dm_build_1378)

	dm_build_1380, dm_build_1381 := dm_build_1376.dm_build_1316(dm_build_1379)
	if dm_build_1381 != nil {
		return nil, dm_build_1381
	}
	return dm_build_1380.(*execInfo), nil
}

func (dm_build_1383 *dm_build_1277) Dm_build_1382(dm_build_1384 *innerRows, dm_build_1385 int64) (*execInfo, error) {
	dm_build_1386 := dm_build_339(dm_build_1383, dm_build_1384, dm_build_1385, INT64_MAX)
	dm_build_1387, dm_build_1388 := dm_build_1383.dm_build_1316(dm_build_1386)
	if dm_build_1388 != nil {
		return nil, dm_build_1388
	}
	return dm_build_1387.(*execInfo), nil
}

func (dm_build_1390 *dm_build_1277) Commit() error {
	dm_build_1391 := dm_build_202(dm_build_1390)
	_, dm_build_1392 := dm_build_1390.dm_build_1316(dm_build_1391)
	if dm_build_1392 != nil {
		return dm_build_1392
	}

	return nil
}

func (dm_build_1394 *dm_build_1277) Rollback() error {
	dm_build_1395 := dm_build_479(dm_build_1394)
	_, dm_build_1396 := dm_build_1394.dm_build_1316(dm_build_1395)
	if dm_build_1396 != nil {
		return dm_build_1396
	}

	return nil
}

func (dm_build_1398 *dm_build_1277) Dm_build_1397(dm_build_1399 *DmConnection) error {
	dm_build_1400 := dm_build_484(dm_build_1398, dm_build_1399.IsoLevel)
	_, dm_build_1401 := dm_build_1398.dm_build_1316(dm_build_1400)
	if dm_build_1401 != nil {
		return dm_build_1401
	}

	return nil
}

func (dm_build_1403 *dm_build_1277) Dm_build_1402(dm_build_1404 *DmStatement, dm_build_1405 string) error {
	dm_build_1406 := dm_build_207(dm_build_1403, dm_build_1404, dm_build_1405)
	_, dm_build_1407 := dm_build_1403.dm_build_1316(dm_build_1406)
	if dm_build_1407 != nil {
		return dm_build_1407
	}

	return nil
}

func (dm_build_1409 *dm_build_1277) Dm_build_1408(dm_build_1410 []uint32) ([]int64, error) {
	dm_build_1411 := dm_build_575(dm_build_1409, dm_build_1410)
	dm_build_1412, dm_build_1413 := dm_build_1409.dm_build_1316(dm_build_1411)
	if dm_build_1413 != nil {
		return nil, dm_build_1413
	}
	return dm_build_1412.([]int64), nil
}

func (dm_build_1415 *dm_build_1277) Close() error {
	if dm_build_1415.dm_build_1287 {
		return nil
	}

	dm_build_1416 := dm_build_1415.dm_build_1278.Close()
	if dm_build_1416 != nil {
		return dm_build_1416
	}

	dm_build_1415.dm_build_1281 = nil
	dm_build_1415.dm_build_1287 = true
	return nil
}

func (dm_build_1418 *dm_build_1277) dm_build_1417(dm_build_1419 *lob) (int64, error) {
	dm_build_1420 := dm_build_370(dm_build_1418, dm_build_1419)
	dm_build_1421, dm_build_1422 := dm_build_1418.dm_build_1316(dm_build_1420)
	if dm_build_1422 != nil {
		return 0, dm_build_1422
	}
	return dm_build_1421.(int64), nil
}

func (dm_build_1424 *dm_build_1277) dm_build_1423(dm_build_1425 *DmBlob, dm_build_1426 int32, dm_build_1427 int32) ([]byte, error) {

	dm_build_1428 := dm_build_357(dm_build_1424, &dm_build_1425.lob, int(dm_build_1426), int(dm_build_1427))
	dm_build_1429, dm_build_1430 := dm_build_1424.dm_build_1316(dm_build_1428)
	if dm_build_1430 != nil {
		return nil, dm_build_1430
	}
	return dm_build_1429.([]byte), nil
}

func (dm_build_1432 *dm_build_1277) dm_build_1431(dm_build_1433 *DmClob, dm_build_1434 int32, dm_build_1435 int32) (string, error) {

	dm_build_1436 := dm_build_357(dm_build_1432, &dm_build_1433.lob, int(dm_build_1434), int(dm_build_1435))
	dm_build_1437, dm_build_1438 := dm_build_1432.dm_build_1316(dm_build_1436)
	if dm_build_1438 != nil {
		return "", dm_build_1438
	}
	dm_build_1439 := dm_build_1437.([]byte)
	return util.Dm_build_586.Dm_build_741(dm_build_1439, 0, len(dm_build_1439), dm_build_1433.serverEncoding), nil
}

func (dm_build_1441 *dm_build_1277) dm_build_1440(dm_build_1442 *DmClob, dm_build_1443 int, dm_build_1444 string, dm_build_1445 string) (int, error) {
	var dm_build_1446 = util.Dm_build_586.Dm_build_793(dm_build_1444, dm_build_1445)
	var dm_build_1447 = 0
	var dm_build_1448 = len(dm_build_1446)
	var dm_build_1449 = 0
	var dm_build_1450 = 0
	var dm_build_1451 = 0
	var dm_build_1452 = dm_build_1448/Dm_build_117 + 1
	var dm_build_1453 byte = 0
	var dm_build_1454 byte = 0x01
	var dm_build_1455 byte = 0x02
	for i := 0; i < dm_build_1452; i++ {
		dm_build_1453 = 0
		if i == 0 {
			dm_build_1453 |= dm_build_1454
		}
		if i == dm_build_1452-1 {
			dm_build_1453 |= dm_build_1455
		}
		dm_build_1451 = dm_build_1448 - dm_build_1450
		if dm_build_1451 > Dm_build_117 {
			dm_build_1451 = Dm_build_117
		}

		setLobData := dm_build_498(dm_build_1441, &dm_build_1442.lob, dm_build_1453, dm_build_1443, dm_build_1446, dm_build_1447, dm_build_1451)
		ret, err := dm_build_1441.dm_build_1316(setLobData)
		if err != nil {
			return 0, err
		}
		tmp := ret.(int32)
		if err != nil {
			return -1, err
		}
		if tmp <= 0 {
			return dm_build_1449, nil
		} else {
			dm_build_1443 += int(tmp)
			dm_build_1449 += int(tmp)
			dm_build_1450 += dm_build_1451
			dm_build_1447 += dm_build_1451
		}
	}
	return dm_build_1449, nil
}

func (dm_build_1457 *dm_build_1277) dm_build_1456(dm_build_1458 *DmBlob, dm_build_1459 int, dm_build_1460 []byte) (int, error) {
	var dm_build_1461 = 0
	var dm_build_1462 = len(dm_build_1460)
	var dm_build_1463 = 0
	var dm_build_1464 = 0
	var dm_build_1465 = 0
	var dm_build_1466 = dm_build_1462/Dm_build_117 + 1
	var dm_build_1467 byte = 0
	var dm_build_1468 byte = 0x01
	var dm_build_1469 byte = 0x02
	for i := 0; i < dm_build_1466; i++ {
		dm_build_1467 = 0
		if i == 0 {
			dm_build_1467 |= dm_build_1468
		}
		if i == dm_build_1466-1 {
			dm_build_1467 |= dm_build_1469
		}
		dm_build_1465 = dm_build_1462 - dm_build_1464
		if dm_build_1465 > Dm_build_117 {
			dm_build_1465 = Dm_build_117
		}

		setLobData := dm_build_498(dm_build_1457, &dm_build_1458.lob, dm_build_1467, dm_build_1459, dm_build_1460, dm_build_1461, dm_build_1465)
		ret, err := dm_build_1457.dm_build_1316(setLobData)
		if err != nil {
			return 0, err
		}
		tmp := ret.(int32)
		if tmp <= 0 {
			return dm_build_1463, nil
		} else {
			dm_build_1459 += int(tmp)
			dm_build_1463 += int(tmp)
			dm_build_1464 += dm_build_1465
			dm_build_1461 += dm_build_1465
		}
	}
	return dm_build_1463, nil
}

func (dm_build_1471 *dm_build_1277) dm_build_1470(dm_build_1472 []interface{}, dm_build_1473 []interface{}, dm_build_1474 int) bool {
	var dm_build_1475 = false
	if dm_build_1473[dm_build_1474] == nil {
		dm_build_1472[dm_build_1474] = ParamDataEnum_Null
	} else if binder, ok := dm_build_1473[dm_build_1474].(iOffRowBinder); ok {
		dm_build_1475 = true
		dm_build_1472[dm_build_1474] = ParamDataEnum_OFF_ROW
		var lob lob
		if l, ok := binder.getObj().(DmBlob); ok {
			lob = l.lob
		} else if l, ok := binder.getObj().(DmClob); ok {
			lob = l.lob
		}
		if &lob != nil && lob.canOptimized(dm_build_1471.dm_build_1281) {
			dm_build_1472[dm_build_1474] = &lobCtl{lob.buildCtlData()}
			dm_build_1475 = false
		}
	} else {
		dm_build_1472[dm_build_1474] = dm_build_1473[dm_build_1474]
	}
	return dm_build_1475
}

func (dm_build_1477 *dm_build_1277) dm_build_1476(dm_build_1478 *DmStatement, dm_build_1479 parameter, dm_build_1480 int, dm_build_1481 iOffRowBinder) error {
	var dm_build_1482 = util.Dm_build_838()
	dm_build_1481.read(dm_build_1482)
	var dm_build_1483 = 0
	for !dm_build_1481.isReadOver() || dm_build_1482.Dm_build_839() > 0 {
		if !dm_build_1481.isReadOver() && dm_build_1482.Dm_build_839() < Dm_build_117 {
			dm_build_1481.read(dm_build_1482)
		}
		if dm_build_1482.Dm_build_839() > Dm_build_117 {
			dm_build_1483 = Dm_build_117
		} else {
			dm_build_1483 = dm_build_1482.Dm_build_839()
		}

		putData := dm_build_469(dm_build_1477, dm_build_1478, int16(dm_build_1480), dm_build_1482, int32(dm_build_1483))
		_, err := dm_build_1477.dm_build_1316(putData)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dm_build_1485 *dm_build_1277) dm_build_1484() ([]byte, error) {
	var dm_build_1486 error
	if dm_build_1485.dm_build_1285 == nil {
		if dm_build_1485.dm_build_1285, dm_build_1486 = security.NewClientKeyPair(); dm_build_1486 != nil {
			return nil, dm_build_1486
		}
	}
	return security.Bn2Bytes(dm_build_1485.dm_build_1285.GetY(), security.DH_KEY_LENGTH), nil
}

func (dm_build_1488 *dm_build_1277) dm_build_1487() (*security.DhKey, error) {
	var dm_build_1489 error
	if dm_build_1488.dm_build_1285 == nil {
		if dm_build_1488.dm_build_1285, dm_build_1489 = security.NewClientKeyPair(); dm_build_1489 != nil {
			return nil, dm_build_1489
		}
	}
	return dm_build_1488.dm_build_1285, nil
}

func (dm_build_1491 *dm_build_1277) dm_build_1490(dm_build_1492 int, dm_build_1493 []byte, dm_build_1494 string, dm_build_1495 int) (dm_build_1496 error) {
	if dm_build_1492 > 0 && dm_build_1492 < security.MIN_EXTERNAL_CIPHER_ID && dm_build_1493 != nil {
		dm_build_1491.dm_build_1282, dm_build_1496 = security.NewSymmCipher(dm_build_1492, dm_build_1493)
	} else if dm_build_1492 >= security.MIN_EXTERNAL_CIPHER_ID {
		if dm_build_1491.dm_build_1282, dm_build_1496 = security.NewThirdPartCipher(dm_build_1492, dm_build_1493, dm_build_1494, dm_build_1495); dm_build_1496 != nil {
			dm_build_1496 = THIRD_PART_CIPHER_INIT_FAILED.addDetailln(dm_build_1496.Error()).throw()
		}
	}
	return
}

func (dm_build_1498 *dm_build_1277) dm_build_1497(dm_build_1499 bool) (dm_build_1500 error) {
	if dm_build_1498.dm_build_1279, dm_build_1500 = security.NewTLSFromTCP(dm_build_1498.dm_build_1278, dm_build_1498.dm_build_1281.dmConnector.sslCertPath, dm_build_1498.dm_build_1281.dmConnector.sslKeyPath, dm_build_1498.dm_build_1281.dmConnector.user); dm_build_1500 != nil {
		return
	}
	if !dm_build_1499 {
		dm_build_1498.dm_build_1279 = nil
	}
	return
}

func (dm_build_1502 *dm_build_1277) dm_build_1501(dm_build_1503 dm_build_123) bool {
	return dm_build_1503.dm_build_137() != Dm_build_33 && dm_build_1502.dm_build_1281.sslEncrypt == 1
}
